package counter

type Manager struct {
	s Store
}

func NewManager(r Store) *Manager {
	return &Manager{s: r}
}

func (m *Manager) Add(id string) error {
	if err := m.s.Store(Counter{ID: id}); err != nil {
		return err
	}

	return nil
}

func (m *Manager) Get(id string) (Counter, error) {
	return m.s.Load(id)
}

func (m *Manager) Inc(id string) error {
	counter, err := m.s.Load(id)
	if err != nil {
		return err
	}

	counter.Inc()

	return m.s.Store(counter)
}

func (m *Manager) Delete(id string) error {
	return m.s.Delete(id)
}
