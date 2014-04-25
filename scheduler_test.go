package tardis

import (
	"gopkg.in/check.v1"
	"time"
)

type SchedulerSuite struct{}
var _ = check.Suite(&SchedulerSuite{})

var (
	scheduler = &Scheduler{Period: 1 * time.Second, Set: Set{Key:"scheduler"}}
)

func (s *SchedulerSuite) SetUpSuite(c *check.C) {
	var err error
	conn := &RedisConn{Address: ":6379"}

	scheduler.Conn, err = conn.Conn()

	if err != nil {
		panic("err connecting to redis on :6379")
	}
}

func (s *SchedulerSuite) SetUpTest(c *check.C) {
	scheduler.Conn.Do("FLUSHALL")
	scheduler.TrackingKey = "tardis:sets"
}

func (s *SchedulerSuite) TestScheduler(c *check.C) {
	
	var ranJob string
	var ranTime int64

	scheduler.Execute = func(job string, timestamp int64) error {
		ranJob = job
		ranTime = timestamp
		return nil
	}
	expectedTime := time.Now().Unix()
	err := scheduler.Add("a job", expectedTime)
	c.Assert(err, check.IsNil)

	scheduler.Start()

   	time.Sleep(2 * time.Second)
	c.Assert(ranJob, check.Equals, "a job")
	c.Assert(ranTime, check.Equals, expectedTime)
}

