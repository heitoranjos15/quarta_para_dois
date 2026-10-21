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
			PlayID:               1,
			GameID:               "2024_01_LAC_DEN",
			Quarter:              1,
			Down:                 1,
			YardsToGo:            10,
			PlayType:             "pass",
			Description:          "Herbert pass complete to Williams for 15 yards",
			EPA:                  f64(1.2),
			Posteam:              "LAC",
			Defteam:              "DEN",
			Yardline100:          75,
			YardsGained:          15,
			PassAttempt:          1,
			RushAttempt:          0,
			CompletePass:         1,
			Touchdown:            0,
			Interception:         0,
			FumbleLost:           0,
			Sack:                 0,
			ThirdDownConverted:   0,
			FourthDownConverted:  0,
			GameSecondsRemaining: 3600,
			Drive:                1,
			HomeTeam:             "LAC",
			AwayTeam:             "DEN",
		},
		{
			PlayID:               2,
			GameID:               "2024_01_LAC_DEN",
			Quarter:              1,
			Down:                 1,
			YardsToGo:            10,
			PlayType:             "run",
			Description:          "Ekeler rush for 5 yards",
			EPA:                  f64(0.3),
			Posteam:              "LAC",
			Defteam:              "DEN",
			Yardline100:          60,
			YardsGained:          5,
			PassAttempt:          0,
			RushAttempt:          1,
			CompletePass:         0,
			Touchdown:            0,
			Interception:         0,
			FumbleLost:           0,
			Sack:                 0,
			ThirdDownConverted:   0,
			FourthDownConverted:  0,
			GameSecondsRemaining: 3560,
			Drive:                1,
			HomeTeam:             "LAC",
			AwayTeam:             "DEN",
		},
		{
			PlayID:               3,
			GameID:               "2024_01_LAC_DEN",
			Quarter:              1,
			Down:                 3,
			YardsToGo:            5,
			PlayType:             "pass",
			Description:          "Herbert pass complete to Allen for 8 yards, 1st down",
			EPA:                  f64(0.8),
			Posteam:              "LAC",
			Defteam:              "DEN",
			Yardline100:          55,
			YardsGained:          8,
			PassAttempt:          1,
			RushAttempt:          0,
			CompletePass:         1,
			Touchdown:            0,
			Interception:         0,
			FumbleLost:           0,
			Sack:                 0,
			ThirdDownConverted:   1,
			FourthDownConverted:  0,
			GameSecondsRemaining: 3520,
			Drive:                1,
			HomeTeam:             "LAC",
			AwayTeam:             "DEN",
		},
		{
			PlayID:               4,
			GameID:               "2024_01_LAC_DEN",
			Quarter:              1,
			Down:                 1,
			YardsToGo:            10,
			PlayType:             "pass",
			Description:          "Wilson pass incomplete",
			EPA:                  f64(-0.5),
			Posteam:              "DEN",
			Defteam:              "LAC",
			Yardline100:          25,
			YardsGained:          0,
			PassAttempt:          1,
			RushAttempt:          0,
			CompletePass:         0,
			Touchdown:            0,
			Interception:         0,
			FumbleLost:           0,
			Sack:                 0,
			ThirdDownConverted:   0,
			FourthDownConverted:  0,
			GameSecondsRemaining: 3480,
			Drive:                1,
			HomeTeam:             "LAC",
			AwayTeam:             "DEN",
		},
		{
			PlayID:               5,
			GameID:               "2024_01_LAC_DEN",
			Quarter:              1,
			Down:                 2,
			YardsToGo:            10,
			PlayType:             "run",
			Description:          "Williams rush for 12 yards",
			EPA:                  f64(1.5),
			Posteam:              "DEN",
			Defteam:              "LAC",
			Yardline100:          25,
			YardsGained:          12,
			PassAttempt:          0,
			RushAttempt:          1,
			CompletePass:         0,
			Touchdown:            0,
			Interception:         0,
			FumbleLost:           0,
			Sack:                 0,
			ThirdDownConverted:   0,
			FourthDownConverted:  0,
			GameSecondsRemaining: 3440,
			Drive:                1,
			HomeTeam:             "LAC",
			AwayTeam:             "DEN",
		},
		{
			PlayID:               6,
			GameID:               "2024_01_LAC_DEN",
			Quarter:              2,
			Down:                 1,
			YardsToGo:            10,
			PlayType:             "pass",
			Description:          "Herbert pass complete to Everett for 25 yards, TD",
			EPA:                  f64(4.2),
			Posteam:              "LAC",
			Defteam:              "DEN",
			Yardline100:          25,
			YardsGained:          25,
			PassAttempt:          1,
			RushAttempt:          0,
			CompletePass:         1,
			Touchdown:            1,
			Interception:         0,
			FumbleLost:           0,
			Sack:                 0,
			ThirdDownConverted:   0,
			FourthDownConverted:  0,
			GameSecondsRemaining: 1800,
			Drive:                2,
			HomeTeam:             "LAC",
			AwayTeam:             "DEN",
		},
		{
			PlayID:               7,
			GameID:               "2024_01_LAC_DEN",
			Quarter:              2,
			Down:                 1,
			YardsToGo:            10,
			PlayType:             "pass",
			Description:          "Wilson intercepted by James",
			EPA:                  f64(-3.0),
			Posteam:              "DEN",
			Defteam:              "LAC",
			Yardline100:          50,
			YardsGained:          0,
			PassAttempt:          1,
			RushAttempt:          0,
			CompletePass:         0,
			Touchdown:            0,
			Interception:         1,
			FumbleLost:           0,
			Sack:                 0,
			ThirdDownConverted:   0,
			FourthDownConverted:  0,
			GameSecondsRemaining: 1760,
			Drive:                2,
			HomeTeam:             "LAC",
			AwayTeam:             "DEN",
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

func TestComputeTeamStats_ThirdDown(t *testing.T) {
	plays := []nflverse.Play{
		{
			PlayID:               1,
			GameID:               "2024_01_LAC_DEN",
			Quarter:              1,
			Down:                 3,
			YardsToGo:            5,
			PlayType:             "pass",
			Description:          "Pass for 6 yards, 1st down",
			EPA:                  f64(1.0),
			Posteam:              "LAC",
			Defteam:              "DEN",
			Yardline100:          50,
			YardsGained:          6,
			PassAttempt:          1,
			RushAttempt:          0,
			CompletePass:         1,
			ThirdDownConverted:   1,
			GameSecondsRemaining: 3600,
			Drive:                1,
			HomeTeam:             "LAC",
			AwayTeam:             "DEN",
		},
		{
			PlayID:               2,
			GameID:               "2024_01_LAC_DEN",
			Quarter:              1,
			Down:                 3,
			YardsToGo:            8,
			PlayType:             "pass",
			Description:          "Pass incomplete",
			EPA:                  f64(-0.8),
			Posteam:              "LAC",
			Defteam:              "DEN",
			Yardline100:          40,
			YardsGained:          0,
			PassAttempt:          1,
			RushAttempt:          0,
			CompletePass:         0,
			ThirdDownConverted:   0,
			GameSecondsRemaining: 3560,
			Drive:                2,
			HomeTeam:             "LAC",
			AwayTeam:             "DEN",
		},
	}

	stats := ComputeTeamStats(plays)

	// 1 of 2 third downs converted = 50%
	assert.InDelta(t, 0.5, stats.HomeOffense.ThirdDownPct, 0.01)
}

func TestComputeTeamStats_ExplosivePlays(t *testing.T) {
	plays := []nflverse.Play{
		{
			PlayID:               1,
			GameID:               "2024_01_LAC_DEN",
			Quarter:              1,
			Down:                 1,
			YardsToGo:            10,
			PlayType:             "pass",
			Description:          "Pass for 25 yards",
			EPA:                  f64(2.0),
			Posteam:              "LAC",
			Defteam:              "DEN",
			Yardline100:          75,
			YardsGained:          25,
			PassAttempt:          1,
			RushAttempt:          0,
			CompletePass:         1,
			ThirdDownConverted:   0,
			GameSecondsRemaining: 3600,
			Drive:                1,
			HomeTeam:             "LAC",
			AwayTeam:             "DEN",
		},
		{
			PlayID:               2,
			GameID:               "2024_01_LAC_DEN",
			Quarter:              1,
			Down:                 1,
			YardsToGo:            10,
			PlayType:             "run",
			Description:          "Run for 15 yards",
			EPA:                  f64(1.5),
			Posteam:              "LAC",
			Defteam:              "DEN",
			Yardline100:          50,
			YardsGained:          15,
			PassAttempt:          0,
			RushAttempt:          1,
			CompletePass:         0,
			ThirdDownConverted:   0,
			GameSecondsRemaining: 3560,
			Drive:                2,
			HomeTeam:             "LAC",
			AwayTeam:             "DEN",
		},
	}

	stats := ComputeTeamStats(plays)

	// 1 explosive pass (20+), 1 explosive rush (10+) = 2 total
	assert.Equal(t, 2, stats.HomeOffense.ExplosivePlays)
}

// offensePlay is a GB-on-offense play in a GB (home) vs DAL (away) game.
func offensePlay(p nflverse.Play) nflverse.Play {
	p.GameID, p.HomeTeam, p.AwayTeam, p.Posteam, p.Defteam = "2026_06_DAL_GB", "GB", "DAL", "GB", "DAL"
	if p.Yardline100 == 0 {
		p.Yardline100 = 60
	}
	if p.FixedDrive == 0 {
		p.FixedDrive = 1
	}
	return p
}

func TestComputeTeamStats_IgnoresNonPlays(t *testing.T) {
	plays := []nflverse.Play{
		offensePlay(nflverse.Play{PlayType: "pass", EPA: f64(1)}),
		offensePlay(nflverse.Play{PlayType: "run", EPA: f64(-1)}),
		offensePlay(nflverse.Play{PlayType: "pass", EPA: f64(1)}),
		offensePlay(nflverse.Play{PlayType: "run", EPA: f64(1)}),
		offensePlay(nflverse.Play{PlayType: "no_play", EPA: f64(5)}),
		offensePlay(nflverse.Play{PlayType: "kickoff", EPA: f64(5)}),
		offensePlay(nflverse.Play{PlayType: "extra_point", EPA: f64(5)}),
		offensePlay(nflverse.Play{PlayType: "timeout", EPA: f64(5)}),
	}

	stats := ComputeTeamStats(plays).HomeOffense

	assert.Equal(t, 4, stats.Plays)
	assert.InDelta(t, 0.5, stats.EPA_per_play, 1e-9)
	assert.InDelta(t, 0.75, stats.SuccessRate, 1e-9)
}

func TestComputeTeamStats_NilEPA(t *testing.T) {
	plays := []nflverse.Play{
		offensePlay(nflverse.Play{PlayType: "pass", EPA: f64(1)}),
		offensePlay(nflverse.Play{PlayType: "run", EPA: f64(-1)}),
		offensePlay(nflverse.Play{PlayType: "pass", EPA: nil}),
	}

	stats := ComputeTeamStats(plays).HomeOffense

	assert.Equal(t, 3, stats.Plays, "a play with no EPA still counts as a play")
	assert.InDelta(t, 2.0/3.0, stats.PassRate, 1e-9)
	assert.InDelta(t, 0.0, stats.EPA_per_play, 1e-9, "averaged over the two plays that have EPA")
	assert.InDelta(t, 0.5, stats.SuccessRate, 1e-9)
}

func TestComputeTeamStats_RedZoneTrips(t *testing.T) {
	plays := []nflverse.Play{
		// drive 1 reaches the red zone and scores a touchdown
		offensePlay(nflverse.Play{PlayType: "pass", FixedDrive: 1, Yardline100: 40, EPA: f64(0.5), FixedDriveResult: "Touchdown"}),
		offensePlay(nflverse.Play{PlayType: "run", FixedDrive: 1, Yardline100: 8, EPA: f64(0.5), FixedDriveResult: "Touchdown"}),
		// drive 2 reaches it twice but only kicks a field goal: still one trip
		offensePlay(nflverse.Play{PlayType: "run", FixedDrive: 2, Yardline100: 18, EPA: f64(0.1), FixedDriveResult: "Field goal"}),
		offensePlay(nflverse.Play{PlayType: "run", FixedDrive: 2, Yardline100: 15, EPA: f64(0.1), FixedDriveResult: "Field goal"}),
		// drive 3 never gets there
		offensePlay(nflverse.Play{PlayType: "pass", FixedDrive: 3, Yardline100: 70, EPA: f64(0), FixedDriveResult: "Punt"}),
	}

	stats := ComputeTeamStats(plays).HomeOffense

	assert.Equal(t, 2, stats.RedZoneTrips)
	assert.Equal(t, 1, stats.RedZoneTDs)
	assert.InDelta(t, 0.5, stats.RedZonePct, 1e-9)
}

func TestComputeTeamStats_NoRedZoneTripsIsZeroPct(t *testing.T) {
	stats := ComputeTeamStats([]nflverse.Play{
		offensePlay(nflverse.Play{PlayType: "pass", EPA: f64(0.2)}),
	}).HomeOffense

	assert.Equal(t, 0, stats.RedZoneTrips)
	assert.Equal(t, 0.0, stats.RedZonePct)
}

func TestComputeTeamStats_TOP(t *testing.T) {
	plays := []nflverse.Play{
		offensePlay(nflverse.Play{PlayType: "pass", FixedDrive: 1, EPA: f64(0), DriveTimeOfPossession: "6:30"}),
		offensePlay(nflverse.Play{PlayType: "run", FixedDrive: 1, EPA: f64(0), DriveTimeOfPossession: "6:30"}),
		offensePlay(nflverse.Play{PlayType: "pass", FixedDrive: 2, EPA: f64(0), DriveTimeOfPossession: "2:07"}),
	}

	stats := ComputeTeamStats(plays).HomeOffense

	assert.Equal(t, 8*time.Minute+37*time.Second, stats.TimeOfPossession)
}
