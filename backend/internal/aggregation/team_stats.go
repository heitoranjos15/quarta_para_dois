package aggregation

import (
	"time"

	"quarta-para-dois/backend/internal/nflverse"
)

type TeamUnitStats struct {
	Plays            int
	EPA_per_play     float64
	SuccessRate      float64
	PassRate         float64
	RushEPA          float64
	PassEPA          float64
	RedZonePct       float64
	ThirdDownPct     float64
	FourthDownPct    float64
	TimeOfPossession time.Duration
	ExplosivePlays   int
	Turnovers        int
}

type TeamComparisonStats struct {
	HomeTeam    string
	AwayTeam    string
	HomeOffense TeamUnitStats
	HomeDefense TeamUnitStats
	AwayOffense TeamUnitStats
	AwayDefense TeamUnitStats
}

func ComputeTeamStats(plays []nflverse.Play) TeamComparisonStats {
	if len(plays) == 0 {
		return TeamComparisonStats{}
	}

	homeTeam := plays[0].HomeTeam
	awayTeam := plays[0].AwayTeam

	stats := TeamComparisonStats{
		HomeTeam: homeTeam,
		AwayTeam: awayTeam,
	}

	homeOffensePlays := filterPlays(plays, func(p nflverse.Play) bool { return p.Posteam == homeTeam })
	homeDefensePlays := filterPlays(plays, func(p nflverse.Play) bool { return p.Defteam == homeTeam })
	awayOffensePlays := filterPlays(plays, func(p nflverse.Play) bool { return p.Posteam == awayTeam })
	awayDefensePlays := filterPlays(plays, func(p nflverse.Play) bool { return p.Defteam == awayTeam })

	stats.HomeOffense = computeUnitStats(homeOffensePlays)
	stats.HomeDefense = computeUnitStats(homeDefensePlays)
	stats.AwayOffense = computeUnitStats(awayOffensePlays)
	stats.AwayDefense = computeUnitStats(awayDefensePlays)

	return stats
}

func filterPlays(plays []nflverse.Play, fn func(nflverse.Play) bool) []nflverse.Play {
	var result []nflverse.Play
	for _, p := range plays {
		if fn(p) {
			result = append(result, p)
		}
	}
	return result
}

func computeUnitStats(plays []nflverse.Play) TeamUnitStats {
	if len(plays) == 0 {
		return TeamUnitStats{}
	}

	stats := TeamUnitStats{
		Plays: len(plays),
	}

	var totalEPA float64
	var successPlays int
	var passPlays int
	var passEPA float64
	var rushEPA float64
	var redZonePlays int
	var redZoneSuccess int
	var thirdDownPlays int
	var thirdDownConverted int
	var fourthDownPlays int
	var fourthDownConverted int
	var lastGameSecondsRemaining map[int]int // drive -> game_seconds_remaining
	var explosivePlays int
	var turnovers int

	lastGameSecondsRemaining = make(map[int]int)

	for _, p := range plays {
		// Total EPA
		totalEPA += p.EPA

		// Success rate (EPA > 0)
		if p.EPA > 0 {
			successPlays++
		}

		// Pass/Rush split
		if p.PassAttempt == 1 {
			passPlays++
			passEPA += p.EPA
		}
		if p.RushAttempt == 1 {
			rushEPA += p.EPA
		}

		// Red zone (yardline100 <= 20)
		if p.Yardline100 <= 20 && p.Yardline100 > 0 {
			redZonePlays++
			if p.EPA > 0 || p.Touchdown == 1 {
				redZoneSuccess++
			}
		}

		// Third down
		if p.Down == 3 {
			thirdDownPlays++
			if p.ThirdDownConverted == 1 {
				thirdDownConverted++
			}
		}

		// Fourth down
		if p.Down == 4 {
			fourthDownPlays++
			if p.FourthDownConverted == 1 {
				fourthDownConverted++
			}
		}

		// Time of possession - track by drive
		if p.Drive > 0 {
			if lastSec, ok := lastGameSecondsRemaining[p.Drive]; ok {
				diff := lastSec - p.GameSecondsRemaining
				if diff > 0 && diff < 300 { // Sanity check: max 5 min per play
					stats.TimeOfPossession += time.Duration(diff) * time.Second
				}
			}
			lastGameSecondsRemaining[p.Drive] = p.GameSecondsRemaining
		}

		// Explosive plays: 20+ pass, 10+ rush
		if p.PassAttempt == 1 && p.YardsGained >= 20 {
			explosivePlays++
		}
		if p.RushAttempt == 1 && p.YardsGained >= 10 {
			explosivePlays++
		}

		// Turnovers
		if p.Interception == 1 || p.FumbleLost == 1 {
			turnovers++
		}
	}

	// Compute averages/percentages
	if stats.Plays > 0 {
		stats.EPA_per_play = totalEPA / float64(stats.Plays)
		stats.SuccessRate = float64(successPlays) / float64(stats.Plays)
		stats.PassRate = float64(passPlays) / float64(stats.Plays)
		stats.PassEPA = passEPA
		stats.RushEPA = rushEPA

		if redZonePlays > 0 {
			stats.RedZonePct = float64(redZoneSuccess) / float64(redZonePlays)
		}
		if thirdDownPlays > 0 {
			stats.ThirdDownPct = float64(thirdDownConverted) / float64(thirdDownPlays)
		}
		if fourthDownPlays > 0 {
			stats.FourthDownPct = float64(fourthDownConverted) / float64(fourthDownPlays)
		}

		stats.ExplosivePlays = explosivePlays
		stats.Turnovers = turnovers
	}

	return stats
}