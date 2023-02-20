package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewHandler(t *testing.T) {
	h := NewHandler()

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
