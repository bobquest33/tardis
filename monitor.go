package tardis

import (
	"github.com/rcrowley/go-metrics"
	"time"
)

type Monitor struct {
	QualifyCount int64
	
	Set
}

func (m *Monitor) Qualify() (bool, error) {
	count, err := m.Set.Count()
	if (err != nil) {
		return false, err
	}
	if count < m.QualifyCount {
		return false, nil
	}
	return true, nil
}

func (m *Monitor) Check(timestamp int64) (int64, int64, error) {
	deltas, err := m.deltas()
	if err != nil {
		return 0, 0, err
	}

	mean := metrics.SampleMean(deltas)
	stdDev := metrics.SampleStdDev(deltas)

	_, _, lastTime, err := m.Set.Last()

	if err != nil {
		return 0, 0, err
	}

	since := timestamp - lastTime - int64(mean)
	
	currentDefcon := int64(float64(since) / stdDev)
	if currentDefcon < 0 {
		currentDefcon = 0
	}
	nextDefconTimestamp := int64(float64(lastTime) + mean + (stdDev * float64(currentDefcon+1)))

	return currentDefcon, nextDefconTimestamp, nil
}

func (m *Monitor) deltas() ([]int64, error) {
	_, times, err := m.Get(0, time.Now().Unix())
	if err != nil {
		return nil, err
	}

	var deltas []int64
	tmp := times[0]

	times = times[1:]
	for _,time := range times {
		deltas = append(deltas, time - tmp)
		tmp = time
	}
	return deltas, nil
}


