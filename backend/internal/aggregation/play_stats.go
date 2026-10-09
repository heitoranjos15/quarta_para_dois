package aggregation

import (
	"sort"
	"time"

	"quarta-para-dois/backend/internal/nflverse"
)

type PlayDetail struct {
	PlayID         int
	Quarter        int
	Down           int
	Distance       int
	Yardline       int
	PlayType       string
	Description    string
	EPA            float64
	Passer         string
	Receiver       string
	Rusher         string
	Tackler        string
	IsScoringPlay  bool
	DriveID        int
	DrivePlayCount int
	DriveYards     int
	DriveResult    string
	DriveTOP       time.Duration
}

type DriveSummary struct {
	DriveID      int
	StartYardline int
	EndYardline   int
	Yards        int
	PlayCount    int
	Result       string
	TOP          time.Duration
}

type QuarterPlays struct {
	Quarter int
	Plays   []PlayDetail
	Drives  []DriveSummary
}

func ProcessPlays(plays []nflverse.Play) []QuarterPlays {
	if len(plays) == 0 {
		return []QuarterPlays{}
	}

	// Sort by play_id
	sortedPlays := make([]nflverse.Play, len(plays))
	copy(sortedPlays, plays)
	sort.Slice(sortedPlays, func(i, j int) bool {
		return sortedPlays[i].PlayID < sortedPlays[j].PlayID
	})

	// Group by quarter
	quarterMap := make(map[int][]nflverse.Play)
	for _, p := range sortedPlays {
		quarterMap[p.Quarter] = append(quarterMap[p.Quarter], p)
	}

	// Get sorted quarters
	var quarters []int
	for q := range quarterMap {
		quarters = append(quarters, q)
	}
	sort.Ints(quarters)

	var result []QuarterPlays
	for _, q := range quarters {
		qPlays := quarterMap[q]
		quarterResult := processQuarter(q, qPlays)
		result = append(result, quarterResult)
	}

	return result
}

func processQuarter(quarter int, plays []nflverse.Play) QuarterPlays {
	var playDetails []PlayDetail
	
	// Group plays by drive
	driveMap := make(map[int][]nflverse.Play)
	for _, p := range plays {
		if p.Drive > 0 {
			driveMap[p.Drive] = append(driveMap[p.Drive], p)
		}
	}

	// Process each play
	for _, p := range plays {
		detail := PlayDetail{
			PlayID:      p.PlayID,
			Quarter:     p.Quarter,
			Down:        p.Down,
			Distance:    p.YardsToGo,
			Yardline:    p.Yardline100,
			PlayType:    p.PlayType,
			Description: p.Description,
			EPA:         p.EPA,
			Passer:      p.PasserPlayerName,
			Receiver:    p.ReceiverPlayerName,
			Rusher:      p.RusherPlayerName,
			Tackler:     p.TacklerPlayerName,
			IsScoringPlay: p.Touchdown == 1,
			DriveID:     p.Drive,
		}
		playDetails = append(playDetails, detail)
	}

	// Compute drive summaries
	var driveSummaries []DriveSummary
	for _, drivePlays := range driveMap {
		drive := computeDriveSummary(drivePlays)
		driveSummaries = append(driveSummaries, drive)
	}

	// Sort drives by drive ID
	sort.Slice(driveSummaries, func(i, j int) bool {
		return driveSummaries[i].DriveID < driveSummaries[j].DriveID
	})

	// Add drive info to play details
	for i := range playDetails {
		for _, ds := range driveSummaries {
			if playDetails[i].DriveID == ds.DriveID {
				playDetails[i].DrivePlayCount = ds.PlayCount
				playDetails[i].DriveYards = ds.Yards
				playDetails[i].DriveResult = ds.Result
				playDetails[i].DriveTOP = ds.TOP
				break
			}
		}
	}

	return QuarterPlays{
		Quarter: quarter,
		Plays:   playDetails,
		Drives:  driveSummaries,
	}
}

func computeDriveSummary(plays []nflverse.Play) DriveSummary {
	if len(plays) == 0 {
		return DriveSummary{}
	}

	// Sort by play_id
	sort.Slice(plays, func(i, j int) bool {
		return plays[i].PlayID < plays[j].PlayID
	})

	drive := DriveSummary{
		DriveID:      plays[0].Drive,
		PlayCount:    len(plays),
		StartYardline: plays[0].Yardline100,
		EndYardline:   plays[len(plays)-1].Yardline100,
	}

	// Calculate yards gained
	for _, p := range plays {
		drive.Yards += p.YardsGained
	}

	// Determine result
	lastPlay := plays[len(plays)-1]
	if lastPlay.Touchdown == 1 {
		drive.Result = "TD"
	} else if lastPlay.Interception == 1 {
		drive.Result = "INT"
	} else if lastPlay.FumbleLost == 1 {
		drive.Result = "FUM"
	} else if lastPlay.Down == 4 && lastPlay.FourthDownConverted == 0 {
		drive.Result = "DOWNS"
	} else if lastPlay.PlayType == "punt" {
		drive.Result = "PUNT"
	} else if lastPlay.PlayType == "field_goal" {
		if lastPlay.YardsGained > 0 { // assuming positive = good
			drive.Result = "FG"
		} else {
			drive.Result = "FG MISS"
		}
	} else if lastPlay.GameSecondsRemaining <= 0 || (lastPlay.Quarter == 4 && lastPlay.GameSecondsRemaining < 120) {
		drive.Result = "END GAME"
	} else if lastPlay.Quarter == 2 && lastPlay.GameSecondsRemaining < 120 {
		drive.Result = "END HALF"
	} else {
		drive.Result = "OTHER"
	}

	// Calculate TOP
	var firstSec, lastSec int
	for _, p := range plays {
		if firstSec == 0 || p.GameSecondsRemaining > firstSec {
			firstSec = p.GameSecondsRemaining
		}
		if lastSec == 0 || p.GameSecondsRemaining < lastSec {
			lastSec = p.GameSecondsRemaining
		}
	}
	if firstSec > lastSec {
		drive.TOP = time.Duration(firstSec-lastSec) * time.Second
	}

	return drive
}