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
	EPA            *float64 // nil when nflverse has no value
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
	DriveID       int
	StartYardline int
	EndYardline   int
	Yards         int
	PlayCount     int
	Result        string
	TOP           time.Duration
}

type QuarterPlays struct {
	Quarter int
	Plays   []PlayDetail
	Drives  []DriveSummary
}

// ProcessPlays lists a game's plays by quarter. Drives are grouped by nflverse's fixed_drive
// across the whole game, so a drive that crosses a quarter break is one drive, listed under
// the quarter it started in.
func ProcessPlays(plays []nflverse.Play) []QuarterPlays {
	if len(plays) == 0 {
		return []QuarterPlays{}
	}

	sortedPlays := make([]nflverse.Play, len(plays))
	copy(sortedPlays, plays)
	sort.Slice(sortedPlays, func(i, j int) bool {
		return sortedPlays[i].PlayID < sortedPlays[j].PlayID
	})

	byDrive := make(map[int][]nflverse.Play)
	for _, p := range sortedPlays {
		if p.FixedDrive > 0 {
			byDrive[p.FixedDrive] = append(byDrive[p.FixedDrive], p)
		}
	}
	lastQuarter := lastGameQuarter(sortedPlays)
	drives := make(map[int]DriveSummary, len(byDrive))
	drivesByQuarter := make(map[int][]DriveSummary)
	for id, drivePlays := range byDrive {
		d := computeDriveSummary(drivePlays, lastQuarter)
		drives[id] = d
		q := driveStartQuarter(drivePlays)
		drivesByQuarter[q] = append(drivesByQuarter[q], d)
	}

	quarterMap := make(map[int][]nflverse.Play)
	for _, p := range sortedPlays {
		quarterMap[p.Quarter] = append(quarterMap[p.Quarter], p)
	}
	var quarters []int
	for q := range quarterMap {
		quarters = append(quarters, q)
	}
	sort.Ints(quarters)

	result := make([]QuarterPlays, 0, len(quarters))
	for _, q := range quarters {
		qDrives := drivesByQuarter[q]
		sort.Slice(qDrives, func(i, j int) bool { return qDrives[i].DriveID < qDrives[j].DriveID })
		result = append(result, QuarterPlays{
			Quarter: q,
			Plays:   playDetails(quarterMap[q], drives),
			Drives:  qDrives,
		})
	}
	return result
}

func playDetails(plays []nflverse.Play, drives map[int]DriveSummary) []PlayDetail {
	out := make([]PlayDetail, 0, len(plays))
	for _, p := range plays {
		d := drives[p.FixedDrive]
		out = append(out, PlayDetail{
			PlayID:         p.PlayID,
			Quarter:        p.Quarter,
			Down:           p.Down,
			Distance:       p.YardsToGo,
			Yardline:       p.Yardline100,
			PlayType:       p.PlayType,
			Description:    p.Description,
			EPA:            p.EPA,
			Passer:         p.PasserPlayerName,
			Receiver:       p.ReceiverPlayerName,
			Rusher:         p.RusherPlayerName,
			Tackler:        p.TacklerPlayerName,
			IsScoringPlay:  p.Touchdown == 1,
			DriveID:        p.FixedDrive,
			DrivePlayCount: d.PlayCount,
			DriveYards:     d.Yards,
			DriveResult:    d.Result,
			DriveTOP:       d.TOP,
		})
	}
	return out
}

// computeDriveSummary summarizes one drive; plays must be sorted by play id.
func computeDriveSummary(plays []nflverse.Play, lastQuarter int) DriveSummary {
	drive := DriveSummary{
		DriveID:       plays[0].FixedDrive,
		PlayCount:     len(plays),
		StartYardline: plays[0].Yardline100,
		EndYardline:   plays[len(plays)-1].Yardline100,
	}
	for _, p := range plays {
		drive.Yards += p.YardsGained
	}
	drive.Result = string(driveResult(firstNonEmpty(plays, func(p nflverse.Play) string { return p.FixedDriveResult }), plays, lastQuarter))
	if top := firstNonEmpty(plays, func(p nflverse.Play) string { return p.DriveTimeOfPossession }); top != "" {
		mmss := parseTime(top)
		drive.TOP = time.Duration(mmss[0])*time.Minute + time.Duration(mmss[1])*time.Second
	}
	return drive
}
