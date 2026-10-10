package aggregation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"quarta-para-dois/backend/internal/nflverse"
)

func loadPlayStatsTestPlays() []nflverse.Play {
	return []nflverse.Play{
		{
			PlayID:             1,
			GameID:             "2024_01_LAC_DEN",
			Quarter:            1,
			Down:               1,
			YardsToGo:          10,
			PlayType:           "pass",
			Description:        "Herbert pass complete to Williams for 15 yards",
			EPA:                1.2,
			Posteam:            "LAC",
			Defteam:            "DEN",
			Yardline100:        75,
			YardsGained:        15,
			PassAttempt:        1,
			RushAttempt:        0,
			CompletePass:       1,
			Touchdown:          0,
			Interception:       0,
			FumbleLost:         0,
			Sack:               0,
			PasserPlayerName:   "J.Herbert",
			ReceiverPlayerName: "M.Williams",
			RusherPlayerName:   "",
			TacklerPlayerName:  "P.Surtain",
			ThirdDownConverted: 0,
			FourthDownConverted: 0,
			GameSecondsRemaining: 3600,
			Drive:              1,
			HomeTeam:           "LAC",
			AwayTeam:           "DEN",
		},
		{
			PlayID:             2,
			GameID:             "2024_01_LAC_DEN",
			Quarter:            1,
			Down:               1,
			YardsToGo:          10,
			PlayType:           "run",
			Description:        "Ekeler rush for 5 yards",
			EPA:                0.3,
			Posteam:            "LAC",
			Defteam:            "DEN",
			Yardline100:        60,
			YardsGained:        5,
			PassAttempt:        0,
			RushAttempt:        1,
			CompletePass:       0,
			Touchdown:          0,
			Interception:       0,
			FumbleLost:         0,
			Sack:               0,
			PasserPlayerName:   "",
			ReceiverPlayerName: "",
			RusherPlayerName:   "A.Ekeler",
			TacklerPlayerName:  "J.Simmons",
			ThirdDownConverted: 0,
			FourthDownConverted: 0,
			GameSecondsRemaining: 3560,
			Drive:              1,
			HomeTeam:           "LAC",
			AwayTeam:           "DEN",
		},
		{
			PlayID:             3,
			GameID:             "2024_01_LAC_DEN",
			Quarter:            1,
			Down:               3,
			YardsToGo:          5,
			PlayType:           "pass",
			Description:        "Herbert pass complete to Allen for 8 yards, 1st down",
			EPA:                0.8,
			Posteam:            "LAC",
			Defteam:            "DEN",
			Yardline100:        55,
			YardsGained:        8,
			PassAttempt:        1,
			RushAttempt:        0,
			CompletePass:       1,
			Touchdown:          0,
			Interception:       0,
			FumbleLost:         0,
			Sack:               0,
			PasserPlayerName:   "J.Herbert",
			ReceiverPlayerName: "K.Allen",
			RusherPlayerName:   "",
			TacklerPlayerName:  "",
			ThirdDownConverted: 1,
			FourthDownConverted: 0,
			GameSecondsRemaining: 3520,
			Drive:              1,
			HomeTeam:           "LAC",
			AwayTeam:           "DEN",
		},
		{
			PlayID:             4,
			GameID:             "2024_01_LAC_DEN",
			Quarter:            1,
			Down:               1,
			YardsToGo:          10,
			PlayType:           "pass",
			Description:        "Herbert pass complete to Everett for 25 yards, TD",
			EPA:                4.2,
			Posteam:            "LAC",
			Defteam:            "DEN",
			Yardline100:        25,
			YardsGained:        25,
			PassAttempt:        1,
			RushAttempt:        0,
			CompletePass:       1,
			Touchdown:          1,
			Interception:       0,
			FumbleLost:         0,
			Sack:               0,
			PasserPlayerName:   "J.Herbert",
			ReceiverPlayerName: "G.Everett",
			RusherPlayerName:   "",
			TacklerPlayerName:  "",
			ThirdDownConverted: 0,
			FourthDownConverted: 0,
			GameSecondsRemaining: 3480,
			Drive:              1,
			HomeTeam:           "LAC",
			AwayTeam:           "DEN",
		},
		{
			PlayID:             5,
			GameID:             "2024_01_LAC_DEN",
			Quarter:            2,
			Down:               1,
			YardsToGo:          10,
			PlayType:           "pass",
			Description:        "Wilson pass incomplete",
			EPA:                -0.5,
			Posteam:            "DEN",
			Defteam:            "LAC",
			Yardline100:        25,
			YardsGained:        0,
			PassAttempt:        1,
			RushAttempt:        0,
			CompletePass:       0,
			Touchdown:          0,
			Interception:       0,
			FumbleLost:         0,
			Sack:               0,
			PasserPlayerName:   "R.Wilson",
			ReceiverPlayerName: "",
			RusherPlayerName:   "",
			TacklerPlayerName:  "D.James",
			ThirdDownConverted: 0,
			FourthDownConverted: 0,
			GameSecondsRemaining: 1800,
			Drive:              2,
			HomeTeam:           "LAC",
			AwayTeam:           "DEN",
		},
		{
			PlayID:             6,
			GameID:             "2024_01_LAC_DEN",
			Quarter:            2,
			Down:               2,
			YardsToGo:          10,
			PlayType:           "run",
			Description:        "Williams rush for 12 yards",
			EPA:                1.5,
			Posteam:            "DEN",
			Defteam:            "LAC",
			Yardline100:        25,
			YardsGained:        12,
			PassAttempt:        0,
			RushAttempt:        1,
			CompletePass:       0,
			Touchdown:          0,
			Interception:       0,
			FumbleLost:         0,
			Sack:               0,
			PasserPlayerName:   "",
			ReceiverPlayerName: "",
			RusherPlayerName:   "J.Williams",
			TacklerPlayerName:  "K.Mack",
			ThirdDownConverted: 0,
			FourthDownConverted: 0,
			GameSecondsRemaining: 1760,
			Drive:              2,
			HomeTeam:           "LAC",
			AwayTeam:           "DEN",
		},
		{
			PlayID:             7,
			GameID:             "2024_01_LAC_DEN",
			Quarter:            2,
			Down:               1,
			YardsToGo:          10,
			PlayType:           "pass",
			Description:        "Wilson intercepted by James",
			EPA:                -3.0,
			Posteam:            "DEN",
			Defteam:            "LAC",
			Yardline100:        50,
			YardsGained:        0,
			PassAttempt:        1,
			RushAttempt:        0,
			CompletePass:       0,
			Touchdown:          0,
			Interception:       1,
			FumbleLost:         0,
			Sack:               0,
			PasserPlayerName:   "R.Wilson",
			ReceiverPlayerName: "",
			RusherPlayerName:   "",
			TacklerPlayerName:  "D.James",
			ThirdDownConverted: 0,
			FourthDownConverted: 0,
			GameSecondsRemaining: 1720,
			Drive:              2,
			HomeTeam:           "LAC",
			AwayTeam:           "DEN",
		},
	}
}

func TestProcessPlays(t *testing.T) {
	plays := loadPlayStatsTestPlays()
	details := ProcessPlays(plays)
	
	require.Len(t, details, 2) // Q1 and Q2
	
	// Check Q1
	q1 := details[0]
	assert.Equal(t, 1, q1.Quarter)
	assert.Len(t, q1.Plays, 4)
	assert.Contains(t, q1.Plays[0].Description, "pass")
	
	// Check Q2
	q2 := details[1]
	assert.Equal(t, 2, q2.Quarter)
	assert.Len(t, q2.Plays, 3)
}

func TestProcessPlays_PlayDetailFields(t *testing.T) {
	plays := loadPlayStatsTestPlays()
	details := ProcessPlays(plays)
	
	q1 := details[0]
	play := q1.Plays[0]
	
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

func TestProcessPlays_ScoringPlay(t *testing.T) {
	plays := loadPlayStatsTestPlays()
	details := ProcessPlays(plays)
	
	// Q1, 4th play is TD
	tdPlay := details[0].Plays[3]
	assert.True(t, tdPlay.IsScoringPlay)
	assert.Equal(t, 4.2, tdPlay.EPA)
}

func TestProcessPlays_Interception(t *testing.T) {
	plays := loadPlayStatsTestPlays()
	details := ProcessPlays(plays)
	
	// Q2, 3rd play is INT
	intPlay := details[1].Plays[2]
	assert.Equal(t, "pass", intPlay.PlayType)
	assert.Equal(t, -3.0, intPlay.EPA)
	assert.Equal(t, "R.Wilson", intPlay.Passer)
}

func TestProcessPlays_DriveGrouping(t *testing.T) {
	plays := loadPlayStatsTestPlays()
	details := ProcessPlays(plays)
	
	// Should have drives in Q1
	q1 := details[0]
	assert.Len(t, q1.Drives, 1)
	
	drive := q1.Drives[0]
	assert.Equal(t, 1, drive.DriveID)
	assert.Equal(t, 4, drive.PlayCount) // 4 plays in drive 1
	assert.Equal(t, 53, drive.Yards) // 15+5+8+25
	assert.Contains(t, drive.Result, "TD") // Ended in TD
	assert.Greater(t, drive.TOP, time.Duration(0))
}

func TestProcessPlays_EmptyPlays(t *testing.T) {
	plays := []nflverse.Play{}
	details := ProcessPlays(plays)
	
	assert.Len(t, details, 0)
}

func TestProcessPlays_SortedByPlayID(t *testing.T) {
	plays := loadPlayStatsTestPlays()
	// Shuffle
	plays[0], plays[1] = plays[1], plays[0]
	plays[2], plays[3] = plays[3], plays[2]
	
	details := ProcessPlays(plays)
	
	// Should be sorted by play_id within quarter
	q1 := details[0]
	for i := 1; i < len(q1.Plays); i++ {
		assert.Less(t, q1.Plays[i-1].PlayID, q1.Plays[i].PlayID, "plays should be sorted by play_id")
	}
}