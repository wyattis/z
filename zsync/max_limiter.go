package zsync

import (
	"fmt"
	"sync"
)

// Keep number of active locks below the Max count. Useful for limiting the
// number of goroutines running simultaneously
func NewMaxLimiter(max int) *MaxLimiter {
	return &MaxLimiter{
		Max:  max,
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

type MaxLimiter struct {
	Active  int
	Pending int
	Max     int
	cond    *sync.Cond
}

// Lock will wait until Active < Max
func (l *MaxLimiter) Lock() {
	l.cond.L.Lock()
	defer l.cond.L.Unlock()
	l.Pending++
	for l.Active >= l.Max {
		if l.Pending >= 2*l.Max {
			fmt.Printf("potential deadlock detected: %+v\n", l)
		}
		l.cond.Wait()
	}
	l.Pending--
	l.Active++
}

// Must call unlock when work is completed
func (l *MaxLimiter) Unlock() {
	l.cond.L.Lock()
	l.Active--
	l.cond.L.Unlock()
	l.cond.Signal()
}
