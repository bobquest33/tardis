package tardis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

type Set struct {
	Key  string
	Conn redis.Conn
	TrackingKey string
}

func (s *Set) Add(value string, timestamp int64) error {
	_, err := s.Conn.Do("ZADD", s.Key, timestamp, value)
	if err != nil {
		return err
	}

	return s.trackLowest()
}

func (s *Set) Remove(value string) error {
	_, err := s.Conn.Do("ZREM", s.Key, value)
	if err != nil {
		return err
	}

	return s.trackLowest()
}

func (s *Set) AddValue(value string) error {
	return s.Add(value, time.Now().Unix())
}

func (s *Set) First() (bool, string, int64, error) {
	return s.GetRank(0)
}

func (s *Set) Last() (bool, string, int64, error) {
	return s.GetRank(-1)
}

func (s *Set) Count() (int64, error) {
	return redis.Int64(s.Conn.Do("ZCARD", s.Key))
}

func (s *Set) GetRank(rank int64) (bool, string, int64, error) {
	lowest, err := redis.Strings(s.Conn.Do("ZRANGE", s.Key, rank, rank, "WITHSCORES"))
	if err != nil {
		return false, "", 0, err
	}
	if len(lowest) == 0 {
		return false, "", 0, nil
	}
	if len(lowest) != 2 {
		return false, "", 0, fmt.Errorf("Unexpected return from ZRANGE: %v", lowest)
	}
	score, err := strconv.ParseInt(lowest[1], 0, 64)
	if err != nil {
		return false, "", 0, err
	}

	return true, lowest[0], score, nil
}

func (s *Set) Get(start_time int64, end_time int64) ([]string, []int64, error) {
	return s.parseResponse(redis.Strings(s.Conn.Do("ZRANGEBYSCORE", s.Key, start_time, end_time, "WITHSCORES")))
}

func (s *Set) GetN(start_time int64, limit int64) ([]string, []int64, error) {
	return s.parseResponse(redis.Strings(s.Conn.Do("ZRANGEBYSCORE",s.Key, "-inf", start_time,"WITHSCORES","LIMIT","0",limit)))
}

func (s *Set) Expire(timestamp int64, fn func(key string, value string, score int64) error) error {
	// race condition here :)
	members, scores, err := s.Get(0, timestamp)
	if err != nil {
		return err
	}

	_, err = redis.Int64(s.Conn.Do("ZREMRANGEBYSCORE", s.Key, 0, timestamp))
	if err != nil {
		return err
	}

	if fn != nil {
		for i, _ := range members {
			err = fn(s.Key, members[i], scores[i])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Clean(trackingKey string, timestamp int64, conn redis.Conn, fn func(key string, value string, score int64) error) error {
	sets := &Set{Key: trackingKey, Conn: conn}

	return sets.Expire(timestamp, func(key string, value string, score int64) error {
		set := &Set{Key: value, Conn: sets.Conn, TrackingKey: trackingKey}
		err := set.Expire(timestamp, fn)
		if err != nil {
			return err
		}
		return set.trackLowest()
	})
}

func (s *Set) trackLowest() error {
	if s.TrackingKey == "" {
		return nil
	}

	exist, _, lowest, err := s.First()
	if err != nil {
		return err
	}

	if exist == false {
		_, err := s.Conn.Do("ZREM", s.TrackingKey, s.Key)
		if err != nil {
			return err
		}
		return nil
	}

	_, err = s.Conn.Do("ZADD", s.TrackingKey, lowest, s.Key)

	if err != nil {
		return err
	}
	return nil
}

func (s *Set) parseResponse(response []string, err error) ([]string, []int64, error) {
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
