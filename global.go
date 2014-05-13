package tardis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

func Epoch(event_time string) (int64, error) {
	// "2014-04-03T20:39:54+00:00"
	parsed, err := time.Parse("2006-01-02T15:04:05+00:00", event_time)
	if err != nil {
		parsed, err = time.Parse("2006-01-02T15:04:05Z", event_time)
		if err != nil {
			return 0, err
		}
	}
	return parsed.Unix(), nil
}

func ParseEvent(event []byte) (map[string]string, error) {
	var raw map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader(event))
	dec.UseNumber()
	err := dec.Decode(&raw)
	if err != nil {
		return nil, err
	}

	parsed := make(map[string]string)

	for k, v := range raw {
		parsed[k] = fmt.Sprintf("%v", v)
	}
	return parsed, nil
}
