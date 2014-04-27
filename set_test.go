package tardis

import (
	"gopkg.in/check.v1"
	"time"
)

var (
	conn RedisConn
	set  = Set{Key: "set1", TrackingKey: "tardis:sets"}
)

type SetSuite struct{}
var _ = check.Suite(&SetSuite{})

func (s *SetSuite) SetUpSuite(c *check.C) {
	var err error
	conn := &RedisConn{Address: ":6379"}

	set.Conn, err = conn.Conn()

	if err != nil {
		panic("err connecting to redis on :6379")
	}
}

func (s *SetSuite) SetUpTest(c *check.C) {
	set.Conn.Do("FLUSHALL")
	set.TrackingKey = "tardis:sets"
}

func (s *SetSuite) TestAddToSet(c *check.C) {
	err := set.Add("1234", 0)
	c.Assert(err, check.IsNil)

	samples, times, err := set.Get(0, 1000)
	c.Assert(err, check.IsNil)
	c.Assert(len(samples), check.Equals, 1)
	c.Assert(len(times), check.Equals, 1)
	c.Assert(samples[0], check.Equals, "1234")
	c.Assert(times[0], check.Equals, int64(0))
}

func (s *SetSuite) TestClean(c *check.C) {
	err := set.Add("1234", 1000)
	c.Assert(err, check.IsNil)

	err = set.Add("5678", 2000)
	c.Assert(err, check.IsNil)

	count, err := set.Count()
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, int64(2))

	err = Clean("tardis:sets", 1500, set.Conn, nil)
	c.Assert(err, check.IsNil)

	samples, times, err := set.Get(0, 4000)
	c.Assert(err, check.IsNil)

	c.Assert(len(samples), check.Equals, 1)
	c.Assert(len(times), check.Equals, 1)
	c.Assert(samples[0], check.Equals, "5678")
	c.Assert(times[0], check.Equals, int64(2000))

	err = Clean("tardis:sets", 3000, set.Conn, nil)
	c.Assert(err, check.IsNil)

	samples, times, err = set.Get(0, 4000)
	c.Assert(err, check.IsNil)

	c.Assert(len(samples), check.Equals, 0)
	c.Assert(len(times), check.Equals, 0)
}

func (s *SetSuite) TestSchedulerPattern(c *check.C) {
	ran := ""

	set.TrackingKey = ""
	set.Key = "scheduler"

	err := set.Add("run body", time.Now().Unix()-5)
	c.Assert(err, check.IsNil)

	err = set.Expire(time.Now().Unix(), func (key string, value string, score int64) error {
		ran = value
		return nil
	})
	c.Assert(err, check.IsNil)
	c.Assert(ran, check.Equals, "run body")
}

func (s *SetSuite) TestRemove(c *check.C) {
	err := set.Add("1234", 1000)
	c.Assert(err, check.IsNil)

	err = set.Add("5678", 2000)
	c.Assert(err, check.IsNil)

	err = set.Remove("1234")
	c.Assert(err, check.IsNil)

	count, err := set.Count()
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, int64(1))

	exist, val, timestamp, err := set.First()
	c.Assert(exist, check.Equals, true)
	c.Assert(err, check.IsNil)
	c.Assert(val, check.Equals, "5678")
	c.Assert(timestamp, check.Equals, int64(2000))

}
