package tardis

import (
  "fmt"
  "time"
)

type Scheduler struct {
  Key     string
  Conn    *RedisConn
  Period  time.Duration
  Execute func(job string, timestamp int64) error
}

func (s *Scheduler) Start() chan bool {
  shutdown := make(chan bool, 1)
  go s.run(shutdown)
  return shutdown
}

func (s *Scheduler) run(shutdown chan bool) {
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
      go s.runJobs()
    }
  }
}

func (s *Scheduler) runJobs() error {
  conn, err := s.Conn.Conn()
  if err != nil {
    return err
  }
  set := &Set{Key: s.Key, Conn: conn, TrackingKey: fmt.Sprintf("scheduler-%s-tracking", s.Key)}
  defer conn.Close()
  err = set.Expire(time.Now().Unix(), func(set string, key string, value int64) error {
    return s.Execute(key, value)
  })
  return err
}

func (s *Scheduler) Add(value string, score int64) error {
  conn, err := s.Conn.Conn()
  if err != nil {
    return err
  }
  set := &Set{Key: s.Key, Conn: conn, TrackingKey: fmt.Sprintf("scheduler-%s-tracking", s.Key)}
  defer conn.Close()
  return set.Add(value, score)
}
