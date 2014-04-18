package tardis

import (
	"testing"
)

var (
	conn Conn
	set  = Set{Key: "set1"}
)

func init() {
	var err error
	conn.Address = ":6379"
	set.Conn, err = conn.Conn()

	if err != nil {
		panic("err connecting to redis")
	}
}

func TestMonitorAddsSampleToSet(t *testing.T) {
	err := set.AddMember("1234", 0)
	if err != nil {
		t.Fatalf("Error while adding sample: %v", err)
	}

	samples, times, err := set.Samples(0, 1000)
	if err != nil {
		t.Fatalf("Error while returning samples: %v", err)
	}

	if len(samples) != 1 || samples[0] != "1234" || len(times) != 1 || times[0] != 0 {
		t.Fatalf("Samples returned incorrect: samples %v, times %v", samples, times)
	}
}

func TestMonitorCleansSamples(t *testing.T) {
	err := set.AddMember("1234", 1000)
	if err != nil {
		t.Fatalf("Error while adding sample: %v", err)
	}

	err = set.AddMember("5678", 2000)
	if err != nil {
		t.Fatalf("Error while adding sample: %v", err)
	}

	Clean(1500, set.Conn)

	samples, times, err := set.Samples(0, 3000)
	if err != nil {
		t.Fatalf("Error while returning samples: %v", err)
	}

	if len(samples) != 1 || len(times) != 1 || samples[0] != "5678" || times[0] != 2000 {
		t.Fatalf("Samples returned incorrect: samples %v, times %v", samples, times)
	}
}
