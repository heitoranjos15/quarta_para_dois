package aggregation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"quarta-para-dois/backend/internal/nflverse"
)

func loadTestPlays() []nflverse.Play {
	return []nflverse.Play{
		{
			PlayID:        1,
			GameID:        "2024_01_LAC_DEN",
			Quarter:       1,
			Down:          1,
			YardsToGo:     10,
			PlayType:      "pass",
			Description:   "Herbert pass complete to Williams for 15 yards",
			EPA:           1.2,
			Posteam:       "LAC",
			Defteam:       "DEN",
			Yardline100:   75,
			YardsGained:   15,
			PassAttempt:   1,
			RushAttempt:   0,
			CompletePass:  1,
			Touchdown:     0,
			Interception:  0,
			FumbleLost:    0,
			Sack:          0,
			ThirdDownConverted: 0,
			FourthDownConverted: 0,
			GameSecondsRemaining: 3600,
			Drive:         1,
			HomeTeam:      "LAC",
			AwayTeam:      "DEN",
		},
		{
			PlayID:        2,
			GameID:        "2024_01_LAC_DEN",
			Quarter:       1,
			Down:          1,
			YardsToGo:     10,
			PlayType:      "run",
			Description:   "Ekeler rush for 5 yards",
			EPA:           0.3,
			Posteam:       "LAC",
			Defteam:       "DEN",
			Yardline100:   60,
			YardsGained:   5,
			PassAttempt:   0,
			RushAttempt:   1,
			CompletePass:  0,
			Touchdown:     0,
			Interception:  0,
			FumbleLost:    0,
			Sack:          0,
			ThirdDownConverted: 0,
			FourthDownConverted: 0,
			GameSecondsRemaining: 3560,
			Drive:         1,
			HomeTeam:      "LAC",
			AwayTeam:      "DEN",
		},
		{
			PlayID:        3,
			GameID:        "2024_01_LAC_DEN",
			Quarter:       1,
			Down:          3,
			YardsToGo:     5,
			PlayType:      "pass",
			Description:   "Herbert pass complete to Allen for 8 yards, 1st down",
			EPA:           0.8,
			Posteam:       "LAC",
			Defteam:       "DEN",
			Yardline100:   55,
			YardsGained:   8,
			PassAttempt:   1,
			RushAttempt:   0,
			CompletePass:  1,
			Touchdown:     0,
			Interception:  0,
			FumbleLost:    0,
			Sack:          0,
			ThirdDownConverted: 1,
			FourthDownConverted: 0,
			GameSecondsRemaining: 3520,
			Drive:         1,
			HomeTeam:      "LAC",
			AwayTeam:      "DEN",
		},
		{
			PlayID:        4,
			GameID:        "2024_01_LAC_DEN",
			Quarter:       1,
			Down:          1,
			YardsToGo:     10,
			PlayType:      "pass",
			Description:   "Wilson pass incomplete",
			EPA:           -0.5,
			Posteam:       "DEN",
			Defteam:       "LAC",
			Yardline100:   25,
			YardsGained:   0,
			PassAttempt:   1,
			RushAttempt:   0,
			CompletePass:  0,
			Touchdown:     0,
			Interception:  0,
			FumbleLost:    0,
			Sack:          0,
			ThirdDownConverted: 0,
			FourthDownConverted: 0,
			GameSecondsRemaining: 3480,
			Drive:         1,
			HomeTeam:      "LAC",
			AwayTeam:      "DEN",
		},
		{
			PlayID:        5,
			GameID:        "2024_01_LAC_DEN",
			Quarter:       1,
			Down:          2,
			YardsToGo:     10,
			PlayType:      "run",
			Description:   "Williams rush for 12 yards",
			EPA:           1.5,
			Posteam:       "DEN",
			Defteam:       "LAC",
			Yardline100:   25,
			YardsGained:   12,
			PassAttempt:   0,
			RushAttempt:   1,
			CompletePass:  0,
			Touchdown:     0,
			Interception:  0,
			FumbleLost:    0,
			Sack:          0,
			ThirdDownConverted: 0,
			FourthDownConverted: 0,
			GameSecondsRemaining: 3440,
			Drive:         1,
			HomeTeam:      "LAC",
			AwayTeam:      "DEN",
		},
		{
			PlayID:        6,
			GameID:        "2024_01_LAC_DEN",
			Quarter:       2,
			Down:          1,
			YardsToGo:     10,
			PlayType:      "pass",
			Description:   "Herbert pass complete to Everett for 25 yards, TD",
			EPA:           4.2,
			Posteam:       "LAC",
			Defteam:       "DEN",
			Yardline100:   25,
			YardsGained:   25,
			PassAttempt:   1,
			RushAttempt:   0,
			CompletePass:  1,
			Touchdown:     1,
			Interception:  0,
			FumbleLost:    0,
			Sack:          0,
			ThirdDownConverted: 0,
			FourthDownConverted: 0,
			GameSecondsRemaining: 1800,
			Drive:         2,
			HomeTeam:      "LAC",
			AwayTeam:      "DEN",
		},
		{
			PlayID:        7,
			GameID:        "2024_01_LAC_DEN",
			Quarter:       2,
			Down:          1,
			YardsToGo:     10,
			PlayType:      "pass",
			Description:   "Wilson intercepted by James",
			EPA:           -3.0,
			Posteam:       "DEN",
			Defteam:       "LAC",
			Yardline100:   50,
			YardsGained:   0,
			PassAttempt:   1,
			RushAttempt:   0,
			CompletePass:  0,
			Touchdown:     0,
			Interception:  1,
			FumbleLost:    0,
			Sack:          0,
			ThirdDownConverted: 0,
			FourthDownConverted: 0,
			GameSecondsRemaining: 1760,
			Drive:         2,
			HomeTeam:      "LAC",
			AwayTeam:      "DEN",
		},
	}
}

func TestComputeTeamStats(t *testing.T) {
	plays := loadTestPlays()
	stats := ComputeTeamStats(plays)
	
	require.NotNil(t, stats)
	assert.Equal(t, "LAC", stats.HomeTeam)
	assert.Equal(t, "DEN", stats.AwayTeam)
	
	// Home offense should have positive EPA
	assert.Greater(t, stats.HomeOffense.EPA_per_play, 0.0)
	// Away defense should have positive success rate (forced negative EPA)
	assert.GreaterOrEqual(t, stats.AwayDefense.SuccessRate, 0.0)
	
	// Check plays counted
	assert.Greater(t, stats.HomeOffense.Plays, 0)
	assert.Greater(t, stats.AwayOffense.Plays, 0)
	
	// Check pass/rush rates
	assert.GreaterOrEqual(t, stats.HomeOffense.PassRate, 0.0)
	assert.LessOrEqual(t, stats.HomeOffense.PassRate, 1.0)
	assert.GreaterOrEqual(t, stats.HomeOffense.RushEPA, -10.0)
}

func TestComputeTeamStats_EmptyPlays(t *testing.T) {
	plays := []nflverse.Play{}
	stats := ComputeTeamStats(plays)
	
	assert.NotNil(t, stats)
	assert.Equal(t, 0, stats.HomeOffense.Plays)
	assert.Equal(t, 0, stats.AwayOffense.Plays)
}

func TestComputeTeamStats_RedZone(t *testing.T) {
	plays := []nflverse.Play{
		{
			PlayID:        1,
			GameID:        "2024_01_LAC_DEN",
			Quarter:       1,
			Down:          1,
			YardsToGo:     10,
			PlayType:      "run",
			Description:   "Run for 5 yards",
			EPA:           0.5,
			Posteam:       "LAC",
			Defteam:       "DEN",
			Yardline100:   15, // Red zone
			YardsGained:   5,
			PassAttempt:   0,
			RushAttempt:   1,
			CompletePass:  0,
			Touchdown:     0,
			ThirdDownConverted: 0,
			FourthDownConverted: 0,
			GameSecondsRemaining: 3600,
			Drive:         1,
			HomeTeam:      "LAC",
			AwayTeam:      "DEN",
		},
		{
			PlayID:        2,
			GameID:        "2024_01_LAC_DEN",
			Quarter:       1,
			Down:          2,
			YardsToGo:     5,
			PlayType:      "pass",
			Description:   "Pass for 5 yards, TD",
			EPA:           2.0,
			Posteam:       "LAC",
			Defteam:       "DEN",
			Yardline100:   10, // Red zone
			YardsGained:   10,
			PassAttempt:   1,
			RushAttempt:   0,
			CompletePass:  1,
			Touchdown:     1,
			ThirdDownConverted: 0,
			FourthDownConverted: 0,
			GameSecondsRemaining: 3560,
			Drive:         1,
			HomeTeam:      "LAC",
			AwayTeam:      "DEN",
		},
	}
	
	stats := ComputeTeamStats(plays)
	
	// Should have red zone plays
	assert.Greater(t, stats.HomeOffense.RedZonePct, 0.0)
}

func TestComputeTeamStats_ThirdDown(t *testing.T) {
	plays := []nflverse.Play{
		{
			PlayID:        1,
			GameID:        "2024_01_LAC_DEN",
			Quarter:       1,
			Down:          3,
			YardsToGo:     5,
			PlayType:      "pass",
			Description:   "Pass for 6 yards, 1st down",
			EPA:           1.0,
			Posteam:       "LAC",
			Defteam:       "DEN",
			Yardline100:   50,
			YardsGained:   6,
			PassAttempt:   1,
			RushAttempt:   0,
			CompletePass:  1,
			ThirdDownConverted: 1,
			GameSecondsRemaining: 3600,
			Drive:         1,
			HomeTeam:      "LAC",
			AwayTeam:      "DEN",
		},
		{
			PlayID:        2,
			GameID:        "2024_01_LAC_DEN",
			Quarter:       1,
			Down:          3,
			YardsToGo:     8,
			PlayType:      "pass",
			Description:   "Pass incomplete",
			EPA:           -0.8,
			Posteam:       "LAC",
			Defteam:       "DEN",
			Yardline100:   40,
			YardsGained:   0,
			PassAttempt:   1,
			RushAttempt:   0,
			CompletePass:  0,
			ThirdDownConverted: 0,
			GameSecondsRemaining: 3560,
			Drive:         2,
			HomeTeam:      "LAC",
			AwayTeam:      "DEN",
		},
	}
	
	stats := ComputeTeamStats(plays)
	
	// 1 of 2 third downs converted = 50%
	assert.InDelta(t, 0.5, stats.HomeOffense.ThirdDownPct, 0.01)
}

func TestComputeTeamStats_ExplosivePlays(t *testing.T) {
	plays := []nflverse.Play{
		{
			PlayID:        1,
			GameID:        "2024_01_LAC_DEN",
			Quarter:       1,
			Down:          1,
			YardsToGo:     10,
			PlayType:      "pass",
			Description:   "Pass for 25 yards",
			EPA:           2.0,
			Posteam:       "LAC",
			Defteam:       "DEN",
			Yardline100:   75,
			YardsGained:   25,
			PassAttempt:   1,
			RushAttempt:   0,
			CompletePass:  1,
			ThirdDownConverted: 0,
			GameSecondsRemaining: 3600,
			Drive:         1,
			HomeTeam:      "LAC",
			AwayTeam:      "DEN",
		},
		{
			PlayID:        2,
			GameID:        "2024_01_LAC_DEN",
			Quarter:       1,
			Down:          1,
			YardsToGo:     10,
			PlayType:      "run",
			Description:   "Run for 15 yards",
			EPA:           1.5,
			Posteam:       "LAC",
			Defteam:       "DEN",
			Yardline100:   50,
			YardsGained:   15,
			PassAttempt:   0,
			RushAttempt:   1,
			CompletePass:  0,
			ThirdDownConverted: 0,
			GameSecondsRemaining: 3560,
			Drive:         2,
			HomeTeam:      "LAC",
			AwayTeam:      "DEN",
		},
	}
	
	stats := ComputeTeamStats(plays)
	
	// 1 explosive pass (20+), 1 explosive rush (10+) = 2 total
	assert.Equal(t, 2, stats.HomeOffense.ExplosivePlays)
}

func TestComputeTeamStats_TimeOfPossession(t *testing.T) {
	plays := []nflverse.Play{
		{
			PlayID:        1,
			GameID:        "2024_01_LAC_DEN",
			Quarter:       1,
			Down:          1,
			YardsToGo:     10,
			PlayType:      "run",
			Description:   "Run for 3 yards",
			EPA:           0.1,
			Posteam:       "LAC",
			Defteam:       "DEN",
			Yardline100:   75,
			YardsGained:   3,
			PassAttempt:   0,
			RushAttempt:   1,
			CompletePass:  0,
			ThirdDownConverted: 0,
			GameSecondsRemaining: 3600,
			Drive:         1,
			HomeTeam:      "LAC",
			AwayTeam:      "DEN",
		},
		{
			PlayID:        2,
			GameID:        "2024_01_LAC_DEN",
			Quarter:       1,
			Down:          2,
			YardsToGo:     7,
			PlayType:      "run",
			Description:   "Run for 4 yards",
			EPA:           0.2,
			Posteam:       "LAC",
			Defteam:       "DEN",
			Yardline100:   72,
			YardsGained:   4,
			PassAttempt:   0,
			RushAttempt:   1,
			CompletePass:  0,
			ThirdDownConverted: 0,
			GameSecondsRemaining: 3560,
			Drive:         1,
			HomeTeam:      "LAC",
			AwayTeam:      "DEN",
		},
	}
	
	stats := ComputeTeamStats(plays)
	
	// Should have positive TOP (40 seconds elapsed on LAC drives)
	assert.Greater(t, stats.HomeOffense.TimeOfPossession, time.Duration(0))
	assert.Equal(t, int64(40*time.Second), int64(stats.HomeOffense.TimeOfPossession))
}