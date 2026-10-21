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
	RedZoneTrips     int     // drives that reached the opponent 20
	RedZoneTDs       int     // of those, drives that ended in a touchdown
	RedZonePct       float64 // RedZoneTDs / RedZoneTrips
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

// isScrimmage is true for the plays team rates are computed on; timeouts, kickoffs,
// extra points and no_play penalties are not plays in that sense.
func isScrimmage(p nflverse.Play) bool {
	return p.PlayType == "pass" || p.PlayType == "run"
}

// computeUnitStats takes every play of one unit (so drive-level numbers can see special teams
// plays too) and computes the per-play rates on its scrimmage plays only.
func computeUnitStats(allPlays []nflverse.Play) TeamUnitStats {
	plays := filterPlays(allPlays, isScrimmage)
	if len(plays) == 0 {
		return TeamUnitStats{}
	}

	stats := TeamUnitStats{Plays: len(plays)}

	var totalEPA, passEPA, rushEPA float64
	var epaPlays, successPlays, passPlays int
	var thirdDownPlays, thirdDownConverted, fourthDownPlays, fourthDownConverted int

	for _, p := range plays {
		if p.PlayType == "pass" {
			passPlays++
		}

		// A play without EPA still counts as a play, but not in the EPA averages.
		if p.EPA != nil {
			epaPlays++
			totalEPA += *p.EPA
			if *p.EPA > 0 {
				successPlays++
			}
			if p.PlayType == "pass" {
				passEPA += *p.EPA
			} else {
				rushEPA += *p.EPA
			}
		}

		if p.Down == 3 {
			thirdDownPlays++
			if p.ThirdDownConverted == 1 {
				thirdDownConverted++
			}
		}
		if p.Down == 4 {
			fourthDownPlays++
			if p.FourthDownConverted == 1 {
				fourthDownConverted++
			}
		}

		// Explosive plays: 20+ pass, 10+ rush
		if p.PlayType == "pass" && p.YardsGained >= 20 {
			stats.ExplosivePlays++
		}
		if p.PlayType == "run" && p.YardsGained >= 10 {
			stats.ExplosivePlays++
		}

		if p.Interception == 1 || p.FumbleLost == 1 {
			stats.Turnovers++
		}
	}

	if epaPlays > 0 {
		stats.EPA_per_play = totalEPA / float64(epaPlays)
		stats.SuccessRate = float64(successPlays) / float64(epaPlays)
	}
	stats.PassRate = float64(passPlays) / float64(stats.Plays)
	stats.PassEPA = passEPA
	stats.RushEPA = rushEPA
	if thirdDownPlays > 0 {
		stats.ThirdDownPct = float64(thirdDownConverted) / float64(thirdDownPlays)
	}
	if fourthDownPlays > 0 {
		stats.FourthDownPct = float64(fourthDownConverted) / float64(fourthDownPlays)
	}

	stats.RedZoneTrips, stats.RedZoneTDs = redZoneTrips(allPlays)
	if stats.RedZoneTrips > 0 {
		stats.RedZonePct = float64(stats.RedZoneTDs) / float64(stats.RedZoneTrips)
	}
	stats.TimeOfPossession = timeOfPossession(allPlays)

	return stats
}

// redZoneTrips counts drives with a play inside the 20, and those that ended in a touchdown.
func redZoneTrips(plays []nflverse.Play) (trips, tds int) {
	reached := map[int]bool{}
	result := map[int]string{}
	for _, p := range plays {
		if p.FixedDrive <= 0 {
			continue
		}
		if p.Yardline100 > 0 && p.Yardline100 <= 20 {
			reached[p.FixedDrive] = true
		}
		if p.FixedDriveResult != "" {
			result[p.FixedDrive] = p.FixedDriveResult
		}
	}
	for drive := range reached {
		trips++
		if result[drive] == "Touchdown" {
			tds++
		}
	}
	return trips, tds
}

// timeOfPossession sums nflverse's per-drive time over the distinct drives in plays.
func timeOfPossession(plays []nflverse.Play) time.Duration {
	seen := map[int]bool{}
	var total time.Duration
	for _, p := range plays {
		if p.FixedDrive <= 0 || seen[p.FixedDrive] || p.DriveTimeOfPossession == "" {
			continue
		}
		seen[p.FixedDrive] = true
		mmss := parseTime(p.DriveTimeOfPossession)
		total += time.Duration(mmss[0])*time.Minute + time.Duration(mmss[1])*time.Second
	}
	return total
}
