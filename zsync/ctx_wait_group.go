package zsync

import (
	"context"
	"sync"
)

type ContextWaitGroup struct {
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

// If any error occurs in this group, the entire group will be cancelled. Each routine must properly handle the
// context cancellation.
func NewContextWaitGroup(ctx context.Context) *ContextWaitGroup {
	ctx, cancel := context.WithCancel(ctx)
	return &ContextWaitGroup{ctx: ctx, cancel: cancel}
}

func (g *ContextWaitGroup) Go(fn func(ctx context.Context) error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		err := fn(g.ctx)
		if err != nil {
			g.errOnce.Do(func() {
				g.err = err
				g.cancel()
			})
		}
	}()
}

func (g *ContextWaitGroup) Wait() error {
	g.wg.Wait()
	return g.err
}

func (g *ContextWaitGroup) Cancel() {
	g.cancel()
}
