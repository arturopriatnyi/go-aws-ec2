package counter

type Counter struct {
	ID    string
	Value uint64
}

func (c *Counter) Inc() {
	c.Value++
}
