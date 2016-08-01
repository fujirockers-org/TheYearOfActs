package util

type Counter struct {
	Count int
}

func (c *Counter) Initialize() {
	c.Count = 1
}

func (c *Counter) Get() int {
	i := c.Count
	c.Count++
	return i
}
