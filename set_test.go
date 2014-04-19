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
	set.Conn.Do("FLUSHALL")
}

func TestMonitorAddsSampleToSet(t *testing.T) {
	err := set.AddMember("1234", 0)
	if err != nil {
		t.Fatalf("Error while adding sample: %v", err)
	}

	samples, times, err := set.Members(0, 1000)
	if err != nil {
		t.Fatalf("Error while returning samples: %v", err)
	}

	if len(samples) != 1 || samples[0] != "1234" || len(times) != 1 || times[0] != 0 {
		t.Fatalf("Members returned incorrect: samples %v, times %v", samples, times)
	}
}

func TestMonitorCleansMembers(t *testing.T) {
	err := set.AddMember("1234", 1000)
	if err != nil {
		t.Fatalf("Error while adding sample: %v", err)
	}

	err = set.AddMember("5678", 2000)
	if err != nil {
		t.Fatalf("Error while adding sample: %v", err)
	}

	err = Clean(1500, set.Conn, nil)
	if err != nil {
		t.Fatalf("Error while cleaning sets: %v", err)
	}

	samples, times, err := set.Members(0, 4000)
	if err != nil {
		t.Fatalf("Error while returning samples: %v", err)
	}

	if len(samples) != 1 || len(times) != 1 || samples[0] != "5678" || times[0] != 2000 {
		t.Fatalf("Members returned incorrect: samples %v, times %v", samples, times)
	}

	err = Clean(3000, set.Conn, nil)
	if err != nil {
		t.Fatalf("Error while cleaning sets: %v", err)
	}

	samples, times, err = set.Members(0, 4000)
	if err != nil {
		t.Fatalf("Error while returning samples: %v", err)
	}

	if len(samples) != 0 || len(times) != 0 {
		t.Fatalf("After second Clean, incorrect state: samples %v, times %v", samples, times)
	}
}
