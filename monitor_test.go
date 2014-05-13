package tardis

import (
	"fmt"
	"gopkg.in/check.v1"
	"testing"
)

var (
	monitor = &Monitor{QualifyCount: 5}

	deltas = []int64{9, 2, 5, 4, 12, 7, 8, 11, 9, 3, 7, 4, 12, 5, 4, 10, 9, 6, 9, 4}
)

func TestMonitor(t *testing.T) { check.TestingT(t) }

type MonitorSuite struct{}

var _ = check.Suite(&MonitorSuite{})

func (s *MonitorSuite) SetUpSuite(c *check.C) {
	var err error
	conn := &RedisConn{Address: ":6379"}

	monitor.Conn, err = conn.Conn()

	if err != nil {
		panic("err connecting to redis on :6379")
	}
}

func (s *MonitorSuite) SetUpTest(c *check.C) {
	monitor.Conn.Do("FLUSHALL")

	var cumulative int64
	cumulative = 0

	for _, delta := range deltas {
		cumulative += delta
		err := monitor.Add(fmt.Sprintf("data-%v", cumulative), cumulative)
		if err != nil {
			panic("err connecting to redis on :6379")
		}
	}
}

func (s *MonitorSuite) TestQualify(c *check.C) {
	qualify, err := monitor.Qualify()
	c.Assert(err, check.IsNil)
	c.Assert(qualify, check.Equals, true)

	monitor.QualifyCount = 500
	qualify, err = monitor.Qualify()
	c.Assert(err, check.IsNil)
	c.Assert(qualify, check.Equals, false)
}

// func (s *MonitorSuite) TestDefConAt(c *check.C) {
// 	defcon, nextTimestamp, err := monitor.Check(125)
// 	c.Assert(err, check.IsNil)
// 	c.Assert(defcon, check.Equals, int64(0))
// 	c.Assert(nextTimestamp, check.Equals, int64(149))

// 	defcon, nextTimestamp, err = monitor.Check(150)
// 	c.Assert(err, check.IsNil)
// 	c.Assert(defcon, check.Equals, int64(1))
// 	c.Assert(nextTimestamp, check.Equals, int64(152))

// 	defcon, nextTimestamp, err = monitor.Check(163)
// 	c.Assert(err, check.IsNil)
// 	c.Assert(defcon, check.Equals, int64(5))
// 	c.Assert(nextTimestamp, check.Equals, int64(165))
// }
