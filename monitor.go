package tardis

type Monitor struct {
  QualifyCount int64

  Set
}

func (m *Monitor) Qualify() (bool, error) {
  count, err := m.Set.Count()
  if err != nil {
    return false, err
  }
  if count < m.QualifyCount {
    return false, nil
  }
  return true, nil
}
