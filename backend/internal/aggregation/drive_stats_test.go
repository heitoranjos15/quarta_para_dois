package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"quarta-para-dois/backend/internal/nflverse"
)

func loadDriveTestPlays() []nflverse.Play {
	return []nflverse.Play{
		{
			PlayID:               1,
			GameID:               "2024_01_LAC_DEN",
			Quarter:              1,
			Down:                 1,
			YardsToGo:            10,
			PlayType:             "pass",
			Description:          "Herbert pass complete to Williams for 15 yards",
			EPA:                  1.2,
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
			EPA:                  0.3,
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
			EPA:                  0.8,
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
			EPA:                  4.2,
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
			EPA:                  -0.5,
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
			EPA:                  1.5,
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
			EPA:                  -3.0,
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
			EPA:                  -0.2,
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
	}
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
	assert.Equal(t, 7, q1.Total["home"].Points)
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
	assert.Equal(t, 1.2, play.EPA)
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
	assert.Equal(t, 4.2, tdPlay.EPA)
}

func TestDetermineEndReason_Touchdown(t *testing.T) {
	lastPlay := nflverse.Play{
		Touchdown: 1,
		PlayType:  "pass",
	}
	reason := determineEndReason(lastPlay, []nflverse.Play{})
	assert.Equal(t, DriveEndTouchdown, reason)
}

func TestDetermineEndReason_Interception(t *testing.T) {
	lastPlay := nflverse.Play{
		Interception: 1,
		PlayType:     "pass",
	}
	reason := determineEndReason(lastPlay, []nflverse.Play{})
	assert.Equal(t, DriveEndInterception, reason)
}

func TestDetermineEndReason_FumbleLost(t *testing.T) {
	lastPlay := nflverse.Play{
		FumbleLost: 1,
		PlayType:   "run",
	}
	reason := determineEndReason(lastPlay, []nflverse.Play{})
	assert.Equal(t, DriveEndFumbleLost, reason)
}

func TestDetermineEndReason_Downs(t *testing.T) {
	lastPlay := nflverse.Play{
		Down:                4,
		FourthDownConverted: 0,
		PlayType:            "pass",
	}
	reason := determineEndReason(lastPlay, []nflverse.Play{})
	assert.Equal(t, DriveEndDowns, reason)
}

func TestDetermineEndReason_Punt(t *testing.T) {
	lastPlay := nflverse.Play{
		PlayType: "punt",
	}
	reason := determineEndReason(lastPlay, []nflverse.Play{})
	assert.Equal(t, DriveEndPunt, reason)
}

func TestDetermineEndReason_FieldGoal(t *testing.T) {
	lastPlay := nflverse.Play{
		PlayType: "field_goal",
	}
	reason := determineEndReason(lastPlay, []nflverse.Play{})
	assert.Equal(t, DriveEndFieldGoal, reason)
}

func TestDetermineEndReason_EndOfHalf(t *testing.T) {
	lastPlay := nflverse.Play{
		Quarter:              2,
		GameSecondsRemaining: 30,
		PlayType:             "qb_kneel",
	}
	reason := determineEndReason(lastPlay, []nflverse.Play{})
	assert.Equal(t, DriveEndEndOfHalf, reason)
}

func TestDetermineEndReason_EndOfGame(t *testing.T) {
	lastPlay := nflverse.Play{
		Quarter:              4,
		GameSecondsRemaining: 10,
		PlayType:             "qb_kneel",
	}
	reason := determineEndReason(lastPlay, []nflverse.Play{})
	assert.Equal(t, DriveEndEndOfGame, reason)
}

func TestDrivePoints(t *testing.T) {
	assert.Equal(t, 7, drivePoints("TOUCHDOWN"))
	assert.Equal(t, 3, drivePoints("FIELD_GOAL"))
	assert.Equal(t, 2, drivePoints("SAFETY"))
	assert.Equal(t, 1, drivePoints("EXTRA_POINT"))
	assert.Equal(t, 0, drivePoints("PUNT"))
	assert.Equal(t, 0, drivePoints("INTERCEPTION"))
	assert.Equal(t, 0, drivePoints("OTHER"))
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
