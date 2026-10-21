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
	DriveEndExtraPoint   DriveEndReason = "EXTRA_POINT"
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
	Overtime *QuarterDriveSummary        `json:"overtime,omitempty"`
}

func ComputeDrives(plays []nflverse.Play, homeTeam, awayTeam string) map[int]QuarterDriveSummary {
	if len(plays) == 0 {
		return make(map[int]QuarterDriveSummary)
	}

	sortedPlays := make([]nflverse.Play, len(plays))
	copy(sortedPlays, plays)
	sort.Slice(sortedPlays, func(i, j int) bool {
		return sortedPlays[i].PlayID < sortedPlays[j].PlayID
	})

	quarterMap := make(map[int][]nflverse.Play)
	for _, p := range sortedPlays {
		quarterMap[p.Quarter] = append(quarterMap[p.Quarter], p)
	}

	var quarters []int
	for q := range quarterMap {
		quarters = append(quarters, q)
	}
	sort.Ints(quarters)

	result := make(map[int]QuarterDriveSummary)
	for _, q := range quarters {
		qPlays := quarterMap[q]
		result[q] = processQuarterDrives(q, qPlays, homeTeam, awayTeam)
	}

	return result
}

func processQuarterDrives(quarter int, plays []nflverse.Play, homeTeam, awayTeam string) QuarterDriveSummary {
	driveMap := make(map[int][]nflverse.Play)
	for _, p := range plays {
		if p.Drive > 0 {
			driveMap[p.Drive] = append(driveMap[p.Drive], p)
		}
	}

	var driveDetails []DriveDetail
	for _, drivePlays := range driveMap {
		detail := computeDriveDetail(drivePlays)
		driveDetails = append(driveDetails, detail)
	}

	sort.Slice(driveDetails, func(i, j int) bool {
		return driveDetails[i].DriveID < driveDetails[j].DriveID
	})

	homeTotal := TeamDriveTotal{}
	awayTotal := TeamDriveTotal{}

	for _, dd := range driveDetails {
		switch dd.Team {
		case homeTeam:
			homeTotal.Drives++
			homeTotal.Yards += dd.Yards
			homeTotal.Points += drivePoints(dd.Result)
			homeTotal.TimeOfPossession = addTimeOfPossession(homeTotal.TimeOfPossession, dd.TimeOfPossession)
		case awayTeam:
			awayTotal.Drives++
			awayTotal.Yards += dd.Yards
			awayTotal.Points += drivePoints(dd.Result)
			awayTotal.TimeOfPossession = addTimeOfPossession(awayTotal.TimeOfPossession, dd.TimeOfPossession)
		}
	}

	return QuarterDriveSummary{
		Drives: driveDetails,
		Total: map[string]TeamDriveTotal{
			"home": homeTotal,
			"away": awayTotal,
		},
	}
}

func computeDriveDetail(plays []nflverse.Play) DriveDetail {
	if len(plays) == 0 {
		return DriveDetail{}
	}

	sort.Slice(plays, func(i, j int) bool {
		return plays[i].PlayID < plays[j].PlayID
	})

	firstPlay := plays[0]
	lastPlay := plays[len(plays)-1]

	var playDetails []PlayDetail
	for _, p := range plays {
		detail := PlayDetail{
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
			DriveID:       p.Drive,
		}
		playDetails = append(playDetails, detail)
	}

	var firstSec, lastSec int
	for _, p := range plays {
		if firstSec == 0 || p.GameSecondsRemaining > firstSec {
			firstSec = p.GameSecondsRemaining
		}
		if lastSec == 0 || p.GameSecondsRemaining < lastSec {
			lastSec = p.GameSecondsRemaining
		}
	}

	top := ""
	if firstSec > lastSec {
		diff := firstSec - lastSec
		top = formatDuration(diff)
	}

	totalYards := 0
	for _, p := range plays {
		totalYards += p.YardsGained
	}

	endReason := determineEndReason(lastPlay, plays)

	return DriveDetail{
		DriveID:            plays[0].Drive,
		Team:               plays[0].Posteam,
		StartYardline:      firstPlay.Yardline100,
		EndYardline:        lastPlay.Yardline100,
		PlaysCount:         len(plays),
		Yards:              totalYards,
		TimeOfPossession:   top,
		Result:             string(endReason),
		EndPlayDescription: lastPlay.Description,
		EndPlayType:        lastPlay.PlayType,
		Plays:              playDetails,
	}
}

func determineEndReason(lastPlay nflverse.Play, allPlays []nflverse.Play) DriveEndReason {
	if lastPlay.Touchdown == 1 {
		return DriveEndTouchdown
	}
	if lastPlay.Interception == 1 {
		return DriveEndInterception
	}
	if lastPlay.FumbleLost == 1 {
		return DriveEndFumbleLost
	}
	if lastPlay.Down == 4 && lastPlay.FourthDownConverted == 0 {
		return DriveEndDowns
	}
	if lastPlay.PlayType == "punt" {
		return DriveEndPunt
	}
	if lastPlay.PlayType == "field_goal" {
		return DriveEndFieldGoal
	}
	if lastPlay.PlayType == "extra_point" {
		return DriveEndExtraPoint
	}
	if lastPlay.Sack == 1 && lastPlay.Yardline100 <= 1 && lastPlay.YardsGained < 0 {
		return DriveEndSafety
	}
	if lastPlay.GameSecondsRemaining <= 0 || (lastPlay.Quarter == 4 && lastPlay.GameSecondsRemaining < 120) {
		return DriveEndEndOfGame
	}
	if lastPlay.Quarter == 2 && lastPlay.GameSecondsRemaining < 120 {
		return DriveEndEndOfHalf
	}
	if lastPlay.Quarter == 4 && lastPlay.GameSecondsRemaining > 0 {
		return DriveEndEndOfGame
	}
	return DriveEndOther
}

func drivePoints(result string) int {
	switch result {
	case "TOUCHDOWN":
		return 7
	case "FIELD_GOAL":
		return 3
	case "SAFETY":
		return 2
	case "EXTRA_POINT":
		return 1
	default:
		return 0
	}
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
