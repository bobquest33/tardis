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

func TestAddToSet(t *testing.T) {
	err := set.Add("1234", 0)
	if err != nil {
		t.Fatalf("Error while adding sample: %v", err)
	}

	samples, times, err := set.Get(0, 1000)
	if err != nil {
		t.Fatalf("Error while returning samples: %v", err)
	}

	if len(samples) != 1 || samples[0] != "1234" || len(times) != 1 || times[0] != 0 {
		t.Fatalf("Get returned incorrect: samples %v, times %v", samples, times)
	}
}

func TestClean(t *testing.T) {
	err := set.Add("1234", 1000)
	if err != nil {
		t.Fatalf("Error while adding sample: %v", err)
	}

	err = set.Add("5678", 2000)
	if err != nil {
		t.Fatalf("Error while adding sample: %v", err)
	}

	err = Clean(1500, set.Conn, nil)
	if err != nil {
		t.Fatalf("Error while cleaning sets: %v", err)
	}

	samples, times, err := set.Get(0, 4000)
	if err != nil {
		t.Fatalf("Error while returning samples: %v", err)
	}

	if len(samples) != 1 || len(times) != 1 || samples[0] != "5678" || times[0] != 2000 {
		t.Fatalf("Get returned incorrect: samples %v, times %v", samples, times)
	}

	err = Clean(3000, set.Conn, nil)
	if err != nil {
		t.Fatalf("Error while cleaning sets: %v", err)
	}

	samples, times, err = set.Get(0, 4000)
	if err != nil {
		t.Fatalf("Error while returning samples: %v", err)
	}

	if len(samples) != 0 || len(times) != 0 {
		t.Fatalf("After second Clean, incorrect state: samples %v, times %v", samples, times)
	}
}
