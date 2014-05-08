package tardis

import (
  "gopkg.in/check.v1"
  "time"
)

type ModelSuite struct {
  model Model
  set Set
}
var _ = check.Suite(&ModelSuite{})

func (s *ModelSuite) SetUpSuite(c *check.C) {
  s.model = Model{
    timeWarp: []Point{Point{ X: 0.0, Y: 0.0},Point{ X: 1.0, Y: 1.0}},
    stat: func(prior []int64, cur int64) float64 {
      return 2.0 
    }}
  s.set = Set { Key: "wut", TrackingKey: "deprecated?"}
  conn := &RedisConn{Address: ":6379"}
  var err error
  s.set.Conn, err = conn.Conn()

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
  pr, err := s.model.Probability(&s.set, time.Now().Unix())
  c.Assert(err,check.IsNil)
  c.Check(pr,check.Equals,  2.0)
}