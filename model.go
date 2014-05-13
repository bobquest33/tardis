package tardis

import (
	"errors"
	"sort"
	"time"
)

type Int64Slice []int64

func (a Int64Slice) Len() int           { return len(a) }
func (a Int64Slice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Int64Slice) Less(i, j int) bool { return a[i] < a[j] }

const (
	hoursInWeek = 7 * 24 * time.Hour
)

type Point struct {
	X, Y float64
}

type Model struct {
	TimeWarp []Point
	Stat     func([]int64, int64) (float64, int64)
	//  length Duration
}

func linearInterpolation(p1 Point, p2 Point, x float64) float64 {
	// y = ax + b
	a := (p2.Y - p1.Y) / (p2.X - p1.X)
	b := p1.Y - a*p1.X
	return a*x + b
}
func startOfWeek(timeStamp int64) time.Time {
	dt := time.Unix(timeStamp, 0)
	weekDay := (dt.Weekday() + 6) % 7
	yr, mth, day := dt.Date()
	return time.Date(yr, mth, (day - int(weekDay)), 0, 0, 0, 0, time.UTC)
}
func (m *Model) WarpTime(t int64) int64 {
	startOfWeek := startOfWeek(t)
	dur := time.Unix(t, 0).Sub(startOfWeek)
	realFraction := float64(dur) / float64(hoursInWeek) // true elapsed percentage of week
	i := sort.Search(len(m.TimeWarp), func(i int) bool { return m.TimeWarp[i].X >= realFraction })
	modelFraction := linearInterpolation(m.TimeWarp[i-1], m.TimeWarp[i], realFraction)        // adjusted percentage of week
	return startOfWeek.Add(time.Duration(int64(modelFraction * float64(hoursInWeek)))).Unix() //model timestamp
}

func (m *Model) UnWarpTime(t int64) int64 {

	startOfWeek := startOfWeek(t)
	dur := time.Unix(t, 0).Sub(startOfWeek)
	modelFraction := float64(dur) / float64(hoursInWeek) // true elapsed percentage of week
	i := sort.Search(len(m.TimeWarp), func(i int) bool { return m.TimeWarp[i].Y >= modelFraction })
	p1 := Point{X: m.TimeWarp[i-1].Y, Y: m.TimeWarp[i-1].X}
	p2 := Point{X: m.TimeWarp[i].X, Y: m.TimeWarp[i].X}
	realFraction := linearInterpolation(p1, p2, modelFraction)                               // adjusted percentage of week
	return startOfWeek.Add(time.Duration(int64(realFraction * float64(hoursInWeek)))).Unix() //model timestamp
}

//Expects array of delta values!
func (m *Model) Probability(mon *Monitor, testPoint int64) (float64, int64, error) {
	_, times, err := mon.GetN(testPoint, 250)
	if err != nil {
		return 0.0, 0, err
	}

	if len(times) == 0 {
		return 0.0, 0, errors.New("Insufficient data.")
	}

	sort.Sort(Int64Slice(times))

	transformedTimes := make([]int64, len(times))
	for i, t := range times {
		transformedTimes[i] = m.WarpTime(t)
	}
	l := transformedTimes[len(transformedTimes)-1]
	p, t := m.Stat(Deltas(transformedTimes), m.WarpTime(testPoint)-l)

	return p, m.UnWarpTime(t) + testPoint, nil
}

func Deltas(times []int64) []int64 {
	var deltas []int64
	tmp := times[0]
	times = times[1:]
	for _, time := range times {
		deltas = append(deltas, time-tmp)
		tmp = time
	}
	return deltas
}
