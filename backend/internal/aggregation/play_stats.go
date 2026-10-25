package aggregation

import (
	"sort"
	"time"

	"quarta-para-dois/backend/internal/nflverse"
)

type PlayDetail struct {
	PlayID         int           `json:"play_id"`
	Quarter        int           `json:"quarter"`
	Down           int           `json:"down"`
	Distance       int           `json:"distance"`
	Yardline       int           `json:"yardline"`
	PlayType       string        `json:"play_type"`
	Description    string        `json:"description"`
	EPA            *float64      `json:"epa"` // nil when nflverse has no value
	Passer         string        `json:"passer"`
	Receiver       string        `json:"receiver"`
	Rusher         string        `json:"rusher"`
	Tackler        string        `json:"tackler"`
	IsScoringPlay  bool          `json:"is_scoring_play"`
	DriveID        int           `json:"drive_id"`
	DrivePlayCount int           `json:"drive_play_count"`
	DriveYards     int           `json:"drive_yards"`
	DriveResult    string        `json:"drive_result"`
	DriveTOP       time.Duration `json:"drive_top_ns"`
}

type DriveSummary struct {
	DriveID       int           `json:"drive_id"`
	StartYardline int           `json:"start_yardline"`
	EndYardline   int           `json:"end_yardline"`
	Yards         int           `json:"yards"`
	PlayCount     int           `json:"play_count"`
	Result        string        `json:"result"`
	TOP           time.Duration `json:"top_ns"`
}

type QuarterPlays struct {
	Quarter int            `json:"quarter"`
	Plays   []PlayDetail   `json:"plays"`
	Drives  []DriveSummary `json:"drives"`
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
