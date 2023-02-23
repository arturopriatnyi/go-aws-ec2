package counter

//go:generate mockgen -source=store.go -destination=mock.go -package=counter

import (
	"errors"
	"sync"
)

var (
	ErrExists    = errors.New("counter exists")
	ErrNotFound  = errors.New("counter not found")
	ErrCorrupted = errors.New("counter corrupted")
)

type Store interface {
	Store(counter Counter) error
	Load(id string) (Counter, error)
	Delete(id string) error
}

type MemoryStore struct {
	data *sync.Map
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: &sync.Map{}}
}

func (s *MemoryStore) Store(counter Counter) error {
	if _, ok := s.data.Load(counter.ID); ok {
		return ErrExists
	}

	s.data.Store(counter.ID, counter)

	return nil
}

func (s *MemoryStore) Load(id string) (Counter, error) {
	v, ok := s.data.Load(id)
	if !ok {
		return Counter{}, ErrNotFound
	}

	counter, ok := v.(Counter)
	if !ok {
		return Counter{}, ErrCorrupted
	}

	return counter, nil
}

func (s *MemoryStore) Delete(id string) error {
	if _, ok := s.data.Load(id); !ok {
		return ErrNotFound
	}

	s.data.Delete(id)

	return nil
}
