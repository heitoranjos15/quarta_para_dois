package nflverse

import (
	"regexp"
	"strconv"
)

var gameIDPattern = regexp.MustCompile(`^(\d{4})_(\d{2})_[A-Z]{2,3}_[A-Z]{2,3}$`)

// ParseGameID splits an nflverse game ID such as 2026_06_WAS_SF into season and week.
func ParseGameID(id string) (season, week int, ok bool) {
	m := gameIDPattern.FindStringSubmatch(id)
	if m == nil {
		return 0, 0, false
	}
	season, _ = strconv.Atoi(m[1])
	week, _ = strconv.Atoi(m[2])
	return season, week, true
}
