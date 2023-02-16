package counter

import (
	"reflect"
	"sync"
	"testing"
)

func TestNewMemoryStore(t *testing.T) {
	store := NewMemoryStore()

	wantStore := &MemoryStore{data: &sync.Map{}}
	if !reflect.DeepEqual(store, wantStore) {
		t.Errorf("want: %+v, got: %+v", wantStore, store)
	}
}

func TestMemoryStore_Store(t *testing.T) {
	for name, tt := range map[string]struct {
		s       func() *MemoryStore
		counter Counter
		wantErr error
	}{
		"OK": {
			s: func() *MemoryStore {
				return &MemoryStore{data: &sync.Map{}}
			},
			counter: Counter{ID: "id", Value: 1},
			wantErr: nil,
		},
		"ErrExists": {
			s: func() *MemoryStore {
				data := &sync.Map{}
				data.Store("id", Counter{ID: "id", Value: 1})

				return &MemoryStore{data: data}
			},
			counter: Counter{ID: "id", Value: 1},
			wantErr: ErrExists,
		},
	} {
		t.Run(name, func(t *testing.T) {
			s := tt.s()

			err := s.Store(tt.counter)

			if err != tt.wantErr {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
			if c, ok := s.data.Load(tt.counter.ID); !ok || !reflect.DeepEqual(c, tt.counter) {
				t.Errorf("want: %+v, got: %+v", tt.counter, c)
			}
		})
	}
}

func TestMemoryStore_Load(t *testing.T) {
	for name, tt := range map[string]struct {
		s           func() *MemoryStore
		id          string
		wantCounter Counter
		wantErr     error
	}{
		"OK": {
			s: func() *MemoryStore {
				data := &sync.Map{}
				data.Store("id", Counter{ID: "id", Value: 1})

				return &MemoryStore{data: data}
			},
			id:          "id",
			wantCounter: Counter{ID: "id", Value: 1},
			wantErr:     nil,
		},
		"ErrNotFound": {
			s: func() *MemoryStore {
				return &MemoryStore{data: &sync.Map{}}
			},
			id:          "id",
			wantCounter: Counter{},
			wantErr:     ErrNotFound,
		},
		"ErrCorrupted": {
			s: func() *MemoryStore {
				data := &sync.Map{}
				data.Store("id", "corrupted counter")

				return &MemoryStore{data: data}
			},
			id:          "id",
			wantCounter: Counter{},
			wantErr:     ErrCorrupted,
		},
	} {
		t.Run(name, func(t *testing.T) {
			s := tt.s()

			c, err := s.Load(tt.id)

			if c != tt.wantCounter {
				t.Errorf("want: %+v, got: %+v", tt.wantCounter, c)
			}
			if err != tt.wantErr {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	for name, tt := range map[string]struct {
		s       func() *MemoryStore
		id      string
		wantErr error
	}{
		"OK": {
			s: func() *MemoryStore {
				data := &sync.Map{}
				data.Store("id", Counter{ID: "id", Value: 1})

				return &MemoryStore{data: data}
			},
			id:      "id",
			wantErr: nil,
		},
		"ErrNotFound": {
			s: func() *MemoryStore {
				return &MemoryStore{data: &sync.Map{}}
			},
			id:      "id",
			wantErr: ErrNotFound,
		},
	} {
		t.Run(name, func(t *testing.T) {
			s := tt.s()

			err := s.Delete(tt.id)

			if err != tt.wantErr {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
			if _, ok := s.data.Load(tt.id); ok {
				t.Error("want: false, got: true")
			}
		})
	}
}
