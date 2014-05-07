package tardis
import (
  "math"
  "time"
  "sort"
)
const (
  hoursInWeek = 7 * 24 * time.Hour
)

type Point struct {
  X, Y float64
}

type Model struct {
  timeWarp []Point
  Statistic func(float64) float64
  //  length Duration
}

func linearInterpolation(p1 Point, p2 Point, x float64) float64{
  // y = ax + b
  a := (p2.Y - p1.X) / (p2.X - p1.X)
  b := p1.Y - a*p1.X
  return math.Floor(a*x + b) 
}

func (m *Model) WarpTime(t int64) int64 {
  dt := time.Unix(t,0)
  weekDay := (dt.Weekday() - 1 ) % 7  - time.Monday - 1 
  yr, mth, day  := dt.Date()
  _, week := dt.ISOWeek()
  startOfWeek := time.Date(yr,mth,week * (day - int(weekDay)),0,0,0,0,time.UTC)
  dur := dt.Sub(startOfWeek)
  realFraction := float64(dur) / float64(hoursInWeek) // true elapsed percentage of week
  i := sort.Search(len(m.timeWarp), func(i int) bool { return m.timeWarp[i].X >= realFraction})
  modelFraction := linearInterpolation(m.timeWarp[i-1],m.timeWarp[i],realFraction) // adjusted percentage of week
  return startOfWeek.Add(time.Duration(int64(modelFraction * float64(hoursInWeek)))).Unix() //model timestamp
}

func (m *Model) UnWarpTime(t int64) int64 {
  //do it all baaaackwards - is that actually needed?
  return 0
}