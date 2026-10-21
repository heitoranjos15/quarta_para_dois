package aggregation

import (
	"fmt"
	"sort"
	"strings"

	"quarta-para-dois/backend/internal/nflverse"
)

type DriveEndReason string

const (
	DriveEndPunt         DriveEndReason = "PUNT"
	DriveEndDowns        DriveEndReason = "DOWNS"
	DriveEndInterception DriveEndReason = "INTERCEPTION"
	DriveEndFumbleLost   DriveEndReason = "FUMBLE_LOST"
	DriveEndFieldGoal    DriveEndReason = "FIELD_GOAL"
	DriveEndMissedFG     DriveEndReason = "MISSED_FIELD_GOAL"
	DriveEndTouchdown    DriveEndReason = "TOUCHDOWN"
	DriveEndSafety       DriveEndReason = "SAFETY"
	DriveEndEndOfHalf    DriveEndReason = "END_OF_HALF"
	DriveEndEndOfGame    DriveEndReason = "END_OF_GAME"
	DriveEndOther        DriveEndReason = "OTHER"
)

type DriveDetail struct {
	DriveID            int          `json:"drive_id"`
	Team               string       `json:"team"`
	StartYardline      int          `json:"start_yardline"`
	EndYardline        int          `json:"end_yardline"`
	PlaysCount         int          `json:"plays_count"`
	Yards              int          `json:"yards"`
	TimeOfPossession   string       `json:"time_of_possession"`
	Result             string       `json:"result"`
	Points             int          `json:"points"` // scored by the offense
	EndPlayDescription string       `json:"end_play_description"`
	EndPlayType        string       `json:"end_play_type"`
	Plays              []PlayDetail `json:"plays"`
}

type TeamDriveTotal struct {
	Drives           int    `json:"drives"`
	Points           int    `json:"points"`
	TimeOfPossession string `json:"time_of_possession"`
	Yards            int    `json:"yards"`
}

type QuarterDriveSummary struct {
	Drives []DriveDetail             `json:"drives"`
	Total  map[string]TeamDriveTotal `json:"total"`
}

type DriveSummaryResponse struct {
	Quarters map[int]QuarterDriveSummary `json:"quarters"`
}

// ComputeDrives groups a game's plays by nflverse's fixed_drive. A drive is listed under the
// quarter it started in (overtime is quarter 5) and counted once in that quarter's totals.
func ComputeDrives(plays []nflverse.Play, homeTeam, awayTeam string) map[int]QuarterDriveSummary {
	byDrive := make(map[int][]nflverse.Play)
	for _, p := range plays {
		if p.FixedDrive > 0 {
			byDrive[p.FixedDrive] = append(byDrive[p.FixedDrive], p)
		}
	}

	byQuarter := make(map[int][]DriveDetail)
	for _, drivePlays := range byDrive {
		sort.Slice(drivePlays, func(i, j int) bool { return drivePlays[i].PlayID < drivePlays[j].PlayID })
		detail := computeDriveDetail(drivePlays)
		quarter := driveStartQuarter(drivePlays)
		byQuarter[quarter] = append(byQuarter[quarter], detail)
	}

	result := make(map[int]QuarterDriveSummary, len(byQuarter))
	for quarter, details := range byQuarter {
		sort.Slice(details, func(i, j int) bool { return details[i].DriveID < details[j].DriveID })
		result[quarter] = summarizeQuarter(details, homeTeam, awayTeam)
	}
	return result
}

func driveStartQuarter(plays []nflverse.Play) int {
	for _, p := range plays {
		if p.DriveQuarterStart > 0 {
			return p.DriveQuarterStart
		}
	}
	return plays[0].Quarter
}

func summarizeQuarter(details []DriveDetail, homeTeam, awayTeam string) QuarterDriveSummary {
	homeTotal := TeamDriveTotal{}
	awayTotal := TeamDriveTotal{}

	for _, dd := range details {
		switch dd.Team {
		case homeTeam:
			addDrive(&homeTotal, dd)
		case awayTeam:
			addDrive(&awayTotal, dd)
		}
	}

	return QuarterDriveSummary{
		Drives: details,
		Total: map[string]TeamDriveTotal{
			"home": homeTotal,
			"away": awayTotal,
		},
	}
}

func addDrive(total *TeamDriveTotal, dd DriveDetail) {
	total.Drives++
	total.Yards += dd.Yards
	total.Points += dd.Points
	total.TimeOfPossession = addTimeOfPossession(total.TimeOfPossession, dd.TimeOfPossession)
}

// computeDriveDetail describes one drive. plays must be sorted by play id.
func computeDriveDetail(plays []nflverse.Play) DriveDetail {
	if len(plays) == 0 {
		return DriveDetail{}
	}

	firstPlay := plays[0]
	lastPlay := plays[len(plays)-1]

	var playDetails []PlayDetail
	totalYards := 0
	for _, p := range plays {
		totalYards += p.YardsGained
		playDetails = append(playDetails, PlayDetail{
			PlayID:        p.PlayID,
			Quarter:       p.Quarter,
			Down:          p.Down,
			Distance:      p.YardsToGo,
			Yardline:      p.Yardline100,
			PlayType:      p.PlayType,
			Description:   p.Description,
			EPA:           epa(p),
			Passer:        p.PasserPlayerName,
			Receiver:      p.ReceiverPlayerName,
			Rusher:        p.RusherPlayerName,
			Tackler:       p.TacklerPlayerName,
			IsScoringPlay: p.Touchdown == 1,
			DriveID:       p.FixedDrive,
		})
	}

	result := driveResult(firstNonEmpty(plays, func(p nflverse.Play) string { return p.FixedDriveResult }), plays)

	return DriveDetail{
		DriveID:            firstPlay.FixedDrive,
		Team:               firstNonEmpty(plays, func(p nflverse.Play) string { return p.Posteam }),
		StartYardline:      firstPlay.Yardline100,
		EndYardline:        lastPlay.Yardline100,
		PlaysCount:         len(plays),
		Yards:              totalYards,
		TimeOfPossession:   firstNonEmpty(plays, func(p nflverse.Play) string { return p.DriveTimeOfPossession }),
		Result:             string(result),
		Points:             drivePoints(result, plays),
		EndPlayDescription: lastPlay.Description,
		EndPlayType:        lastPlay.PlayType,
		Plays:              playDetails,
	}
}

func firstNonEmpty(plays []nflverse.Play, get func(nflverse.Play) string) string {
	for _, p := range plays {
		if v := get(p); v != "" {
			return v
		}
	}
	return ""
}

// driveResult maps nflverse's fixed_drive_result to a DriveEndReason. plays must be the
// drive's plays in order; they decide what a generic "Turnover" was.
func driveResult(fixedDriveResult string, plays []nflverse.Play) DriveEndReason {
	switch fixedDriveResult {
	case "Touchdown":
		return DriveEndTouchdown
	case "Field goal":
		return DriveEndFieldGoal
	case "Missed field goal":
		return DriveEndMissedFG
	case "Punt":
		return DriveEndPunt
	case "Turnover on downs":
		return DriveEndDowns
	case "Safety":
		return DriveEndSafety
	case "Turnover", "Opp touchdown":
		return turnoverKind(plays)
	case "End of half":
		if len(plays) > 0 && plays[len(plays)-1].DriveQuarterEnd >= 4 {
			return DriveEndEndOfGame
		}
		return DriveEndEndOfHalf
	}
	return DriveEndOther
}

func turnoverKind(plays []nflverse.Play) DriveEndReason {
	for _, p := range plays {
		if p.Interception == 1 {
			return DriveEndInterception
		}
	}
	for _, p := range plays {
		if p.FumbleLost == 1 {
			return DriveEndFumbleLost
		}
	}
	if len(plays) > 0 && plays[len(plays)-1].PlayType == "punt" {
		return DriveEndPunt
	}
	return DriveEndOther
}

// drivePoints is what the offense scored on the drive. A defensive touchdown scores nothing for it.
func drivePoints(result DriveEndReason, plays []nflverse.Play) int {
	switch result {
	case DriveEndTouchdown:
		points := 6
		for _, p := range plays {
			if p.ExtraPointResult == "good" {
				points++
			}
			if p.TwoPointConvResult == "success" {
				points += 2
			}
		}
		return points
	case DriveEndFieldGoal:
		return 3
	}
	return 0
}

func formatDuration(seconds int) string {
	if seconds <= 0 {
		return "0:00"
	}
	minutes := seconds / 60
	secs := seconds % 60
	return formatTimeOfPossession(minutes, secs)
}

func formatTimeOfPossession(minutes, seconds int) string {
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func addTimeOfPossession(current, addition string) string {
	if current == "" {
		return addition
	}
	if addition == "" {
		return current
	}

	cParts := parseTime(current)
	aParts := parseTime(addition)

	totalSecs := cParts[0]*60 + cParts[1] + aParts[0]*60 + aParts[1]
	return formatTimeOfPossession(totalSecs/60, totalSecs%60)
}

func parseTime(t string) []int {
	parts := strings.Split(t, ":")
	if len(parts) != 2 {
		return []int{0, 0}
	}
	minutes := 0
	seconds := 0
	_, _ = fmt.Sscanf(parts[0], "%d", &minutes)
	_, _ = fmt.Sscanf(parts[1], "%d", &seconds)
	return []int{minutes, seconds}
}
