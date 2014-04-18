package tardis

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

type Conn struct {
	Pool    *redis.Pool
	Address string
}

func (c *Conn) Conn() (redis.Conn, error) {
	var conn redis.Conn
	var err error

	if c.Pool != nil {
		conn = c.Pool.Get()
	} else {
		conn, err = redis.Dial("tcp", c.Address)
		if err != nil {
			return nil, err
		}
	}
	return conn, nil
}

func (r *Conn) InitPool(size int) error {
	r.Pool = &redis.Pool{
		MaxIdle:     size,
		IdleTimeout: 30 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", r.Address)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}
