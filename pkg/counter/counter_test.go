package counter

import "testing"

func TestCounter_Inc(t *testing.T) {
	var c Counter

	c.Inc()

	if c.Value != 1 {
		t.Errorf("want: %d, got: %d", 1, c.Value)
	}
}
