package tardis

import (
  "gopkg.in/check.v1"
  "time"
)

type ModelSuite struct {
  model Model
  monitor Monitor
}
var _ = check.Suite(&ModelSuite{})

func (s *ModelSuite) SetUpSuite(c *check.C) {
  s.model = Model{
    TimeWarp: []Point{Point{ X: 0.0, Y: 0.0},Point{ X: 1.0, Y: 1.0}},
    Stat: func(prior []int64, cur int64) (float64, int64) {
      return 2.0, 0 
    }}
  s.monitor = Monitor{Set: Set{ Key: "wut", TrackingKey: "deprecated?"}, QualifyCount: 1}
  conn := &RedisConn{Address: ":6379"}
  var err error
  s.monitor.Set.Conn, err = conn.Conn()

  if err != nil {
    panic("err connecting to redis on :6379")
  }
}

func (s *ModelSuite) TestStartOfWeek(c *check.C) {
  t := startOfWeek(time.Now().Unix())
  day := t.Weekday()
  c.Check(day,check.Equals,time.Monday)
  c.Check(t.Hour(),check.Equals,0)
  c.Check(t.Minute(),check.Equals,0)
  t, _ =time.Parse(time.RFC3339, "2014-04-27T04:19:36Z" )
  t = startOfWeek(t.Unix())
  day = t.Weekday()
  c.Check(day,check.Equals,time.Monday)
  c.Check(t.Day(),check.Equals, 21)
}

func (s *ModelSuite) TestLinearInterpolation(c *check.C) {
  p1 := Point{X: 0.0, Y: 0.0}
  p2 := Point{X: 1.0, Y: 1.0}
  fX := linearInterpolation(p1,p2,0.5)
  c.Check(fX,check.Equals,0.5)
  p1 = Point{X: 1.0, Y: 2.0}
  p2 = Point{X: 7.0, Y: 6.0}
  fX = linearInterpolation(p1,p2,4.0)
  c.Check(fX,check.Equals,4.0)
}

func (s *ModelSuite) TestProbability(c *check.C) {
  _, _, err:= s.model.Probability(&s.monitor,time.Now().Unix())
  c.Assert(err,check.NotNil)

  // c.Check(pr,check.Equals,  2.0)
}