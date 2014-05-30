package tardis

import (
	"fmt"
)

var (
	testDeltas = []int64{9, 2, 5, 4, 12, 7, 8, 11, 9, 3, 7, 4, 12, 5, 4, 10, 9, 6, 9, 4}
)

func insertDeltas(m *Monitor, deltas []int64) {
	var cumulative int64
	cumulative = 0

	for _, delta := range deltas {
		cumulative += delta
		err := m.Add(fmt.Sprintf("data-%v", cumulative), cumulative)
		if err != nil {
			panic("err connecting to redis on :6379")
		}
	}
}
