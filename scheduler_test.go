package tardis

import (
  "gopkg.in/check.v1"
  "time"
)

type SchedulerSuite struct{}

var _ = check.Suite(&SchedulerSuite{})

var (
  scheduler = &Scheduler{Period: 1 * time.Second, Key: "scheduler", Conn: &RedisConn{Address: ":6379"}}
)

func (s *SchedulerSuite) SetUpTest(c *check.C) {
  conn, err := scheduler.Conn.Conn()
  if err != nil {
    panic("Error connecting to redis")
  }
  defer conn.Close()
  conn.Do("FLUSHALL")
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
