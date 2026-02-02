package youtube

import (
	"errors"
	"regexp"
	"strconv"
)

var isoDurationRegex = regexp.MustCompile(
	`^PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?$`,
)

func ParseISODurationToSeconds(iso string) (int, error) {
	matches := isoDurationRegex.FindStringSubmatch(iso)
	if matches == nil {
		return 0, errors.New("invalid ISO-8601 duration")
	}

	toInt := func(s string) int {
		if s == "" {
			return 0
		}
		v, _ := strconv.Atoi(s)
		return v
	}

	hours := toInt(matches[1])
	minutes := toInt(matches[2])
	seconds := toInt(matches[3])

	total := hours*3600 + minutes*60 + seconds
	return total, nil
}
