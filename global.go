package tardis

import (
	"time"
)

func Epoch(event_time string) (int64, error) {
	// "2014-04-03T20:39:54+00:00"
	parsed, err := time.Parse("2006-01-02T15:04:05+00:00", event_time)
	if err != nil {
		return 0, err
	}
	return parsed.Unix(), nil
}
