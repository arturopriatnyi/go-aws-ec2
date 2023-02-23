package counter

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNewManager(t *testing.T) {
	s := NewMockStore(gomock.NewController(t))

	m := NewManager(s)

	if !reflect.DeepEqual(m, &Manager{s: s}) {
		t.Errorf("want: %+v, got: %+v", &Manager{s: s}, m)
	}
}

func TestManager_Add(t *testing.T) {
	for name, tt := range map[string]struct {
		s       func(*gomock.Controller) Store
		id      string
		wantErr error
	}{
		"OK": {
			s: func(c *gomock.Controller) Store {
				s := NewMockStore(c)

				s.
					EXPECT().
					Store(Counter{ID: "id"}).
					Return(nil)

				return s
			},
			id:      "id",
			wantErr: nil,
		},
		"ErrExists": {
			id: "id",
			s: func(c *gomock.Controller) Store {
				s := NewMockStore(c)

				s.
					EXPECT().
					Store(Counter{ID: "id"}).
					Return(ErrExists)

				return s
			},
			wantErr: ErrExists,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := &Manager{s: tt.s(gomock.NewController(t))}

			err := m.Add(tt.id)

			if err != tt.wantErr {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestManager_Get(t *testing.T) {
	for name, tt := range map[string]struct {
		s           func(*gomock.Controller) Store
		id          string
		wantCounter Counter
		wantErr     error
	}{
		"OK": {
			s: func(c *gomock.Controller) Store {
				s := NewMockStore(c)

				s.
					EXPECT().
					Load("id").
					Return(Counter{ID: "id", Value: 1}, nil)

				return s
			},
			id:          "id",
			wantCounter: Counter{ID: "id", Value: 1},
			wantErr:     nil,
		},
		"ErrNotFound": {
			s: func(c *gomock.Controller) Store {
				s := NewMockStore(c)

				s.
					EXPECT().
					Load("id").
					Return(Counter{}, ErrNotFound)

				return s
			},
			id:          "id",
			wantCounter: Counter{},
			wantErr:     ErrNotFound,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := &Manager{s: tt.s(gomock.NewController(t))}

			c, err := m.Get(tt.id)

			if c != tt.wantCounter {
				t.Errorf("want: %+v, got: %+v", tt.wantCounter, c)
			}
			if err != tt.wantErr {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestManager_Inc(t *testing.T) {
	for name, tt := range map[string]struct {
		s       func(*gomock.Controller) Store
		id      string
		wantErr error
	}{
		"OK": {
			s: func(c *gomock.Controller) Store {
				s := NewMockStore(c)

				s.
					EXPECT().
					Load("id").
					Return(Counter{ID: "id", Value: 1}, nil)
				s.
					EXPECT().
					Delete("id").
					Return(nil)
				s.
					EXPECT().
					Store(Counter{ID: "id", Value: 2}).
					Return(nil)

				return s
			},
			id:      "id",
			wantErr: nil,
		},
		"LoadErrNotFound": {
			s: func(c *gomock.Controller) Store {
				s := NewMockStore(c)

				s.
					EXPECT().
					Load("id").
					Return(Counter{}, ErrNotFound)

				return s
			},
			id:      "id",
			wantErr: ErrNotFound,
		},
		"DeleteErrNotFound": {
			s: func(c *gomock.Controller) Store {
				s := NewMockStore(c)

				s.
					EXPECT().
					Load("id").
					Return(Counter{ID: "id", Value: 1}, nil)
				s.
					EXPECT().
					Delete("id").
					Return(ErrNotFound)

				return s
			},
			id:      "id",
			wantErr: ErrNotFound,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := &Manager{s: tt.s(gomock.NewController(t))}

			err := m.Inc(tt.id)

			if err != tt.wantErr {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestManager_Delete(t *testing.T) {
	for name, tt := range map[string]struct {
		s       func(*gomock.Controller) Store
		id      string
		wantErr error
	}{
		"OK": {
			s: func(c *gomock.Controller) Store {
				s := NewMockStore(c)

				s.
					EXPECT().
					Delete("id").
					Return(nil)

				return s
			},
			id:      "id",
			wantErr: nil,
		},
		"ErrNotFound": {
			s: func(c *gomock.Controller) Store {
				s := NewMockStore(c)

				s.
					EXPECT().
					Delete("id").
					Return(ErrNotFound)

				return s
			},
			id:      "id",
			wantErr: ErrNotFound,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := &Manager{s: tt.s(gomock.NewController(t))}

			err := m.Delete(tt.id)

			if err != tt.wantErr {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
