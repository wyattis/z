package zsignal

import (
	"fmt"
	"time"
)

func New() *Signal {
	return &Signal{
		make(chan bool),
	}
}

type Signal struct {
	tickle chan bool
}

func (s *Signal) Notify() {
	s.tickle <- true
}

func (s *Signal) Every(handle func() error) {
	for range s.tickle {
		if err := handle(); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (s *Signal) Debounce(minDelay time.Duration, handle func() error) {
	t := time.NewTicker(minDelay)
	wasTickled := false
	for {
		select {
		case <-t.C:
			if !wasTickled {
				continue
			}
			if err := handle(); err != nil {
				fmt.Println(err)
			}
			wasTickled = false
		case <-s.tickle:
			wasTickled = true
		}
	}
}
