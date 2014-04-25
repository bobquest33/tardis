package tardis

import (
	"time"
)

type Scheduler struct {
	Set
	Period time.Duration
	Execute func(job string, timestamp int64) error
}

func (s *Scheduler) Start() chan bool {
	shutdown := make (chan bool, 1)
	go s.run(shutdown)
	return shutdown
}

func (s *Scheduler) run(shutdown chan bool) {
	conn := s.Conn
	defer conn.Close()
	timeout := make(chan bool, 1)
	
	for true {
		go func() {
    			time.Sleep(s.Period)
    			timeout <- true
		}()

		select {
		case <-shutdown:
			return
		case <-timeout:
			s.runJobs()
		}
	}
}

func (s *Scheduler) runJobs() error {
	err := s.Expire(time.Now().Unix(), func (set string, key string, value int64) error {
		return s.Execute(key, value)
	})
	return err
}

