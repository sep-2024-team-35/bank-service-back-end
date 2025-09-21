package utils

import (
	"fmt"
	"time"
)

func ParseMerchantTimestamp(ts string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05.000Z07:00",
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05",
	}

	var t time.Time
	var err error

	for _, layout := range layouts {
		t, err = time.Parse(layout, ts)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("cannot parse merchant timestamp: %s", ts)
}
