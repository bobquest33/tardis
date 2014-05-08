package tardis
import (
// "math"
"time"
"sort"
// "fmt"
)
const (
  hoursInWeek = 7 * 24 * time.Hour
)

type Point struct {
  X, Y float64
}

type Model struct {
  timeWarp []Point
  stat func([]int64, int64) float64
  //  length Duration
}

func linearInterpolation(p1 Point, p2 Point, x float64) float64{
  // y = ax + b
  a := (p2.Y - p1.Y) / (p2.X - p1.X)
  b := p1.Y - a*p1.X 
  return a*x + b
}
func startOfWeek(timeStamp int64) time.Time {
  dt := time.Unix(timeStamp,0)
  weekDay := (dt.Weekday() - 1 ) % 7
  yr, mth, day  := dt.Date()
  return time.Date(yr,mth, (day - int(weekDay)),0,0,0,0,time.UTC)  
}
func (m *Model) TimeWarp(t int64) int64 {
  startOfWeek := startOfWeek(t)
  dur := time.Unix(t,0).Sub(startOfWeek)
  realFraction := float64(dur) / float64(hoursInWeek) // true elapsed percentage of week
  i := sort.Search(len(m.timeWarp), func(i int) bool { return m.timeWarp[i].X >= realFraction})
  modelFraction := linearInterpolation(m.timeWarp[i-1],m.timeWarp[i],realFraction) // adjusted percentage of week
  return startOfWeek.Add(time.Duration(int64(modelFraction * float64(hoursInWeek)))).Unix() //model timestamp
}

func (m *Model) TimeUnWarp(t int64) int64 {
  //do it all baaaackwards - is that actually needed?
  return 0
}

func(m *Model) Probability(points *Set, testPoint int64) (float64, error){
  _, times, err := points.GetN(testPoint,150)
  if(err != nil) {
    return 0.0, err
  }
  transformedTimes := make([]int64,len(times))
  for i, t := range times {
    transformedTimes[i] = m.TimeWarp(t)
  }
  return m.stat(transformedTimes,m.TimeWarp(testPoint)), nil
}