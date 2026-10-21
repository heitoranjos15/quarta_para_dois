package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"quarta-para-dois/backend/internal/nflverse"
)

func loadDriveTestPlays() []nflverse.Play {
	return withFixedDrives([]nflverse.Play{
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
			PasserPlayerName:     "J.Herbert",
			ReceiverPlayerName:   "M.Williams",
			RusherPlayerName:     "",
			TacklerPlayerName:    "P.Surtain",
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
			PasserPlayerName:     "",
			ReceiverPlayerName:   "",
			RusherPlayerName:     "A.Ekeler",
			TacklerPlayerName:    "J.Simmons",
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
			PasserPlayerName:     "J.Herbert",
			ReceiverPlayerName:   "K.Allen",
			RusherPlayerName:     "",
			TacklerPlayerName:    "",
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
			PasserPlayerName:     "J.Herbert",
			ReceiverPlayerName:   "G.Everett",
			RusherPlayerName:     "",
			TacklerPlayerName:    "",
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
			Quarter:              2,
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
			PasserPlayerName:     "R.Wilson",
			ReceiverPlayerName:   "",
			RusherPlayerName:     "",
			TacklerPlayerName:    "D.James",
			ThirdDownConverted:   0,
			FourthDownConverted:  0,
			GameSecondsRemaining: 1800,
			Drive:                2,
			HomeTeam:             "LAC",
			AwayTeam:             "DEN",
		},
		{
			PlayID:               6,
			GameID:               "2024_01_LAC_DEN",
			Quarter:              2,
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
			PasserPlayerName:     "",
			ReceiverPlayerName:   "",
			RusherPlayerName:     "J.Williams",
			TacklerPlayerName:    "K.Mack",
			ThirdDownConverted:   0,
			FourthDownConverted:  0,
			GameSecondsRemaining: 1760,
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
			PasserPlayerName:     "R.Wilson",
			ReceiverPlayerName:   "",
			RusherPlayerName:     "",
			TacklerPlayerName:    "D.James",
			ThirdDownConverted:   0,
			FourthDownConverted:  0,
			GameSecondsRemaining: 1720,
			Drive:                2,
			HomeTeam:             "LAC",
			AwayTeam:             "DEN",
		},
		{
			PlayID:               8,
			GameID:               "2024_01_LAC_DEN",
			Quarter:              2,
			Down:                 1,
			YardsToGo:            10,
			PlayType:             "punt",
			Description:          "Dixon punts 45 yards",
			EPA:                  f64(-0.2),
			Posteam:              "DEN",
			Defteam:              "LAC",
			Yardline100:          30,
			YardsGained:          0,
			PassAttempt:          0,
			RushAttempt:          0,
			CompletePass:         0,
			Touchdown:            0,
			Interception:         0,
			FumbleLost:           0,
			Sack:                 0,
			PasserPlayerName:     "",
			ReceiverPlayerName:   "",
			RusherPlayerName:     "",
			TacklerPlayerName:    "",
			ThirdDownConverted:   0,
			FourthDownConverted:  0,
			GameSecondsRemaining: 1680,
			Drive:                3,
			HomeTeam:             "LAC",
			AwayTeam:             "DEN",
		},
	})
}

func TestComputeDrives(t *testing.T) {
	plays := loadDriveTestPlays()
	quarters := ComputeDrives(plays, "LAC", "DEN")

	require.Len(t, quarters, 2)

	q1 := quarters[1]
	require.Len(t, q1.Drives, 1)

	drive1 := q1.Drives[0]
	assert.Equal(t, 1, drive1.DriveID)
	assert.Equal(t, "LAC", drive1.Team)
	assert.Equal(t, 75, drive1.StartYardline)
	assert.Equal(t, 25, drive1.EndYardline)
	assert.Equal(t, 4, drive1.PlaysCount)
	assert.Equal(t, 53, drive1.Yards)
	assert.Equal(t, "TOUCHDOWN", drive1.Result)
	assert.Equal(t, "Herbert pass complete to Everett for 25 yards, TD", drive1.EndPlayDescription)
	assert.Equal(t, "pass", drive1.EndPlayType)
	assert.Len(t, drive1.Plays, 4)

	assert.Equal(t, 1, q1.Total["home"].Drives)
	assert.Equal(t, 6, q1.Total["home"].Points) // the fixture has no extra-point play
	assert.Equal(t, 53, q1.Total["home"].Yards)

	q2 := quarters[2]
	require.Len(t, q2.Drives, 2)

	drive2 := q2.Drives[0]
	assert.Equal(t, 2, drive2.DriveID)
	assert.Equal(t, "DEN", drive2.Team)
	assert.Equal(t, "INTERCEPTION", drive2.Result)

	drive3 := q2.Drives[1]
	assert.Equal(t, 3, drive3.DriveID)
	assert.Equal(t, "DEN", drive3.Team)
	assert.Equal(t, "PUNT", drive3.Result)

	assert.Equal(t, 2, q2.Total["away"].Drives)
	assert.Equal(t, 0, q2.Total["away"].Points)
}

func TestComputeDrives_EmptyPlays(t *testing.T) {
	plays := []nflverse.Play{}
	quarters := ComputeDrives(plays, "LAC", "DEN")

	assert.Len(t, quarters, 0)
}

func TestComputeDrives_SortedByDriveID(t *testing.T) {
	plays := loadDriveTestPlays()
	quarters := ComputeDrives(plays, "LAC", "DEN")

	q2 := quarters[2]
	for i := 1; i < len(q2.Drives); i++ {
		assert.Less(t, q2.Drives[i-1].DriveID, q2.Drives[i].DriveID, "drives should be sorted by drive_id")
	}
}

func TestComputeDrives_PlayDetailFields(t *testing.T) {
	plays := loadDriveTestPlays()
	quarters := ComputeDrives(plays, "LAC", "DEN")

	q1 := quarters[1]
	drive := q1.Drives[0]
	play := drive.Plays[0]

	assert.Equal(t, 1, play.PlayID)
	assert.Equal(t, 1, play.Quarter)
	assert.Equal(t, 1, play.Down)
	assert.Equal(t, 10, play.Distance)
	assert.Equal(t, 75, play.Yardline)
	assert.Equal(t, "pass", play.PlayType)
	assert.Contains(t, play.Description, "Herbert")
	require.NotNil(t, play.EPA)
	assert.Equal(t, 1.2, *play.EPA)
	assert.Equal(t, "J.Herbert", play.Passer)
	assert.Equal(t, "M.Williams", play.Receiver)
	assert.Equal(t, "", play.Rusher)
	assert.Equal(t, "P.Surtain", play.Tackler)
	assert.False(t, play.IsScoringPlay)
	assert.Equal(t, 1, play.DriveID)
}

func TestComputeDrives_ScoringPlay(t *testing.T) {
	plays := loadDriveTestPlays()
	quarters := ComputeDrives(plays, "LAC", "DEN")

	q1 := quarters[1]
	drive := q1.Drives[0]
	tdPlay := drive.Plays[3]

	assert.True(t, tdPlay.IsScoringPlay)
	require.NotNil(t, tdPlay.EPA)
	assert.Equal(t, 4.2, *tdPlay.EPA)
}

func TestFormatDuration(t *testing.T) {
	assert.Equal(t, "0:00", formatDuration(0))
	assert.Equal(t, "0:00", formatDuration(-5))
	assert.Equal(t, "0:30", formatDuration(30))
	assert.Equal(t, "1:00", formatDuration(60))
	assert.Equal(t, "1:30", formatDuration(90))
	assert.Equal(t, "3:45", formatDuration(225))
	assert.Equal(t, "10:00", formatDuration(600))
}

func TestAddTimeOfPossession(t *testing.T) {
	assert.Equal(t, "1:30", addTimeOfPossession("1:00", "0:30"))
	assert.Equal(t, "2:45", addTimeOfPossession("1:30", "1:15"))
	assert.Equal(t, "0:45", addTimeOfPossession("", "0:45"))
	assert.Equal(t, "1:00", addTimeOfPossession("1:00", ""))
}

// withFixedDrives fills the nflverse drive columns the way the real file does: fixed_drive
// mirrors drive, and every play of a drive carries the drive's result and starting quarter.
func withFixedDrives(plays []nflverse.Play) []nflverse.Play {
	results := map[int]string{1: "Touchdown", 2: "Turnover", 3: "Punt"}
	times := map[int]string{1: "2:40", 2: "1:15", 3: "0:45"}
	for i := range plays {
		plays[i].FixedDrive = plays[i].Drive
		plays[i].FixedDriveResult = results[plays[i].Drive]
		plays[i].DriveTimeOfPossession = times[plays[i].Drive]
		plays[i].DriveQuarterStart = plays[i].Quarter
		plays[i].DriveQuarterEnd = plays[i].Quarter
	}
	return plays
}

func onePlayDrive(result string, plays ...nflverse.Play) []nflverse.Play {
	for i := range plays {
		plays[i].GameID = "2026_06_DAL_GB"
		plays[i].PlayID = i + 1
		plays[i].FixedDrive = 1
		plays[i].FixedDriveResult = result
		plays[i].Posteam = "GB"
		if plays[i].Quarter == 0 {
			plays[i].Quarter = 1
		}
		if plays[i].DriveQuarterStart == 0 {
			plays[i].DriveQuarterStart = plays[i].Quarter
		}
		if plays[i].DriveQuarterEnd == 0 {
			plays[i].DriveQuarterEnd = plays[i].Quarter
		}
	}
	return plays
}

func onlyDrive(t *testing.T, plays []nflverse.Play) (DriveDetail, QuarterDriveSummary) {
	t.Helper()
	quarters := ComputeDrives(plays, "GB", "DAL")
	for _, q := range quarters {
		if len(q.Drives) == 1 {
			return q.Drives[0], q
		}
	}
	t.Fatalf("expected exactly one drive, got %+v", quarters)
	return DriveDetail{}, QuarterDriveSummary{}
}

func TestComputeDrives_FourthDownPuntIsPunt(t *testing.T) {
	plays := onePlayDrive("Punt",
		nflverse.Play{PlayType: "pass", Down: 3},
		nflverse.Play{PlayType: "punt", Down: 4},
	)

	drive, q := onlyDrive(t, plays)

	assert.Equal(t, "PUNT", drive.Result)
	assert.Equal(t, 0, drive.Points)
	assert.Equal(t, 0, q.Total["home"].Points)
}

func TestComputeDrives_FieldGoalMadeAndMissed(t *testing.T) {
	made, q := onlyDrive(t, onePlayDrive("Field goal",
		nflverse.Play{PlayType: "run", Down: 1},
		nflverse.Play{PlayType: "field_goal", Down: 4, FieldGoalResult: "made"},
	))
	assert.Equal(t, "FIELD_GOAL", made.Result)
	assert.Equal(t, 3, made.Points)
	assert.Equal(t, 3, q.Total["home"].Points)

	missed, _ := onlyDrive(t, onePlayDrive("Missed field goal",
		nflverse.Play{PlayType: "field_goal", Down: 4, FieldGoalResult: "missed"},
	))
	assert.Equal(t, "MISSED_FIELD_GOAL", missed.Result)
	assert.Equal(t, 0, missed.Points)
}

func TestComputeDrives_PickSix(t *testing.T) {
	drive, q := onlyDrive(t, onePlayDrive("Opp touchdown",
		nflverse.Play{PlayType: "pass", Down: 2},
		nflverse.Play{PlayType: "pass", Down: 3, Interception: 1, Touchdown: 1, TDTeam: "DAL"},
	))

	assert.Equal(t, "INTERCEPTION", drive.Result)
	assert.Equal(t, 0, drive.Points, "the defense scored, not the offense")
	assert.Equal(t, 0, q.Total["home"].Points)
}

func TestComputeDrives_TouchdownPoints(t *testing.T) {
	td := nflverse.Play{PlayType: "pass", Down: 2, Touchdown: 1}
	cases := []struct {
		name  string
		extra nflverse.Play
		want  int
	}{
		{"extra point good", nflverse.Play{PlayType: "extra_point", ExtraPointResult: "good"}, 7},
		{"two point success", nflverse.Play{PlayType: "run", TwoPointConvResult: "success"}, 8},
		{"extra point failed", nflverse.Play{PlayType: "extra_point", ExtraPointResult: "failed"}, 6},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			drive, _ := onlyDrive(t, onePlayDrive("Touchdown", td, tc.extra))
			assert.Equal(t, "TOUCHDOWN", drive.Result)
			assert.Equal(t, tc.want, drive.Points)
		})
	}
}

func TestComputeDrives_CrossQuarter(t *testing.T) {
	plays := onePlayDrive("Field goal",
		nflverse.Play{PlayType: "pass", Down: 1, Quarter: 1, GameSecondsRemaining: 3000},
		nflverse.Play{PlayType: "run", Down: 1, Quarter: 2, GameSecondsRemaining: 2700},
		nflverse.Play{PlayType: "field_goal", Down: 4, Quarter: 2, FieldGoalResult: "made", GameSecondsRemaining: 2600},
	)
	for i := range plays {
		plays[i].DriveQuarterStart, plays[i].DriveQuarterEnd = 1, 2
	}

	quarters := ComputeDrives(plays, "GB", "DAL")

	require.Len(t, quarters, 1)
	q1 := quarters[1]
	require.Len(t, q1.Drives, 1)
	assert.Equal(t, 3, q1.Drives[0].PlaysCount)
	assert.Equal(t, 1, q1.Total["home"].Drives)
	assert.NotContains(t, quarters, 2, "the drive is counted once, in the quarter it started")
}

func TestComputeDrives_Overtime(t *testing.T) {
	plays := onePlayDrive("End of half",
		nflverse.Play{PlayType: "pass", Down: 1, Quarter: 5},
		nflverse.Play{PlayType: "run", Down: 2, Quarter: 5},
	)

	quarters := ComputeDrives(plays, "GB", "DAL")

	require.Contains(t, quarters, 5)
	assert.Equal(t, "END_OF_GAME", quarters[5].Drives[0].Result)
}

func TestComputeDrives_EndOfFirstHalf(t *testing.T) {
	drive, _ := onlyDrive(t, onePlayDrive("End of half",
		nflverse.Play{PlayType: "qb_kneel", Down: 1, Quarter: 2},
	))

	assert.Equal(t, "END_OF_HALF", drive.Result)
}

func TestComputeDrives_TimeOfPossessionFromNflverse(t *testing.T) {
	plays := onePlayDrive("Punt", nflverse.Play{PlayType: "punt", Down: 4})
	plays[0].DriveTimeOfPossession = "2:35"

	drive, _ := onlyDrive(t, plays)

	assert.Equal(t, "2:35", drive.TimeOfPossession)
}

func TestComputeDrives_EndOfRegulationBeforeOvertime(t *testing.T) {
	regulation := onePlayDrive("End of half",
		nflverse.Play{PlayType: "qb_kneel", Down: 1, Quarter: 4},
	)
	overtime := onePlayDrive("Field goal",
		nflverse.Play{PlayType: "field_goal", Down: 4, Quarter: 5, FieldGoalResult: "made"},
	)
	overtime[0].FixedDrive, overtime[0].PlayID = 2, 10

	quarters := ComputeDrives(append(regulation, overtime...), "GB", "DAL")

	require.Contains(t, quarters, 4)
	assert.Equal(t, "END_OF_HALF", quarters[4].Drives[0].Result, "the game went on to overtime")
	assert.Equal(t, "FIELD_GOAL", quarters[5].Drives[0].Result)
}

func TestComputeDrives_EndOfRegulationIsEndOfGame(t *testing.T) {
	drive, _ := onlyDrive(t, onePlayDrive("End of half",
		nflverse.Play{PlayType: "qb_kneel", Down: 1, Quarter: 4},
	))

	assert.Equal(t, "END_OF_GAME", drive.Result)
}
