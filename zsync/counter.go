package zsync

import "sync"

func NewCounter() *Counter {
	return &Counter{
		cond: &sync.Cond{
			L: &sync.Mutex{},
		},
	}
}

type Counter struct {
	Value int
	cond  *sync.Cond
}

func (c *Counter) Set(val int) {
	c.cond.L.Lock()
	defer c.cond.L.Unlock()
	c.Value += val
	c.cond.Broadcast()
}

func (c *Counter) Add(val int) {
	c.cond.L.Lock()
	defer c.cond.L.Unlock()
	c.Value += val
	c.cond.Broadcast()
}

func (c *Counter) WaitExactly(val int) <-chan bool {
	out := make(chan bool)
	go func() {
		c.cond.L.Lock()
		defer c.cond.L.Unlock()
		for c.Value != val {
			c.cond.Wait()
		}
		out <- true
		close(out)
	}()
	return out
}
