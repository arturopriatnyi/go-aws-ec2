package http

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-aws-ec2/pkg/counter"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestNewHandler(t *testing.T) {
	h := NewHandler(zap.NewNop(), NewMockCounterManager(gomock.NewController(t)))

	if h == nil {
		t.Errorf("want handler: <non-nil>, got: <nil>")
	}
}

func Test_noRoute(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	noRoute()(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("want status code: 404, got: %d", w.Code)
	}
	if w.Body.String() != "" {
		t.Errorf("want body: , got: %s", w.Body.String())
	}
}

func Test_noMethod(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	noMethod()(c)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("want status code: 405, got: %d", w.Code)
	}
	if w.Body.String() != "" {
		t.Errorf("want body: , got: %s", w.Body.String())
	}
}

func Test_health(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	health()(c)

	if w.Code != 200 {
		t.Errorf("want status code: 200, got: %d", w.Code)
	}
	if w.Body.String() != "" {
		t.Errorf("want body: , got: %s", w.Body.String())
	}
}

func Test_addCounter(t *testing.T) {
	for name, tt := range map[string]struct {
		cm       func(c *gomock.Controller) CounterManager
		body     string
		wantCode int
		wantBody string
	}{
		"OK": {
			cm: func(c *gomock.Controller) CounterManager {
				cm := NewMockCounterManager(c)

				cm.EXPECT().Add("id").Return(nil)

				return cm
			},
			body:     `{"id":"id"}`,
			wantCode: http.StatusCreated,
			wantBody: "",
		},
		"BadRequestInvalidBody": {
			cm: func(c *gomock.Controller) CounterManager {
				return NewMockCounterManager(c)
			},
			body:     ``,
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"EOF"}`,
		},
		"BadRequestCounterExists": {
			cm: func(c *gomock.Controller) CounterManager {
				cm := NewMockCounterManager(c)

				cm.EXPECT().Add("id").Return(counter.ErrExists)

				return cm
			},
			body:     `{"id":"id"}`,
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"counter exists"}`,
		},
		"InternalServerError": {
			cm: func(c *gomock.Controller) CounterManager {
				cm := NewMockCounterManager(c)

				cm.EXPECT().Add("id").Return(errors.New("unexpected error"))

				return cm
			},
			body:     `{"id":"id"}`,
			wantCode: http.StatusInternalServerError,
			wantBody: "",
		},
	} {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = &http.Request{
				Body: io.NopCloser(bytes.NewBufferString(tt.body)),
			}

			addCounter(zap.NewNop(), tt.cm(gomock.NewController(t)))(c)

			if w.Code != tt.wantCode {
				t.Errorf("want status code: %d, got: %d", tt.wantCode, w.Code)
			}
			if w.Body.String() != tt.wantBody {
				t.Errorf("want body: %s, got: %s", tt.wantBody, w.Body.String())
			}
		})
	}
}
