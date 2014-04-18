package tardis

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
        "time"
)

var (
	SetKey = "tardis:sets"
)

type Set struct {
	Key  string
	Conn redis.Conn
}

func (s *Set) AddMember(value string, timestamp int64) error {
	_, err := s.Conn.Do("ZADD", s.Key, timestamp, value)
	if err != nil {
		return err
	}

	_, err = s.Conn.Do("ZADD", SetKey, timestamp, s.Key)
	if err != nil {
		return err
	}
	return nil
}

func (s *Set) AddValue(value string) error {
	return s.AddMember(value, time.Now().Unix())
}

func (s *Set) Samples(start_time int64, end_time int64) ([]string, []int64, error) {
	response, err := redis.Strings(s.Conn.Do("ZRANGEBYSCORE", s.Key, start_time, end_time, "WITHSCORES"))
	if err != nil {
		return nil, nil, err
	}
	var events []string
	var times []int64
	for i, resp := range response {
		if i%2 == 0 {
			events = append(events, resp)
		} else {
			val, err := strconv.ParseInt(resp, 0, 64)
			if err != nil {
				return nil, nil, err
			}
			times = append(times, val)
		}
	}

	return events, times, nil
}

func Clean(timestamp int64, conn redis.Conn) error {

	sets, err := redis.Strings(conn.Do("ZRANGEBYSCORE", SetKey, 0, timestamp))
	if err != nil {
		return err
	}

	for _, set := range sets {
		conn.Do("ZREMRANGEBYSCORE", set, 0, timestamp)
		if err != nil {
			return err
		}
	}
	conn.Do("ZREMRANGEBYSCORE", SetKey, 0, timestamp)
	if err != nil {
		return err
	}
	return nil
}
