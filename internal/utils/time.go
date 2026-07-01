package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}

	d, err := time.ParseDuration(s)
	if err == nil {
		return d, nil
	}

	// otherwise we handle d, w, m, y suffixes

	suffixes := map[string]time.Duration{
		"d": 24 * time.Hour,
		"w": 7 * 24 * time.Hour,
		"m": 30 * 24 * time.Hour,
		"y": 365 * 24 * time.Hour,
	}

	for suff, dur := range suffixes {
		if before, ok :=strings.CutSuffix(s, suff); ok  {
			numStr := before
			num, err := strconv.Atoi(numStr)
			if err != nil {
				continue
			}
			return time.Duration(num) * dur, nil
		}
	}
	return 0, fmt.Errorf("invalid duration format: %s", s)
}
