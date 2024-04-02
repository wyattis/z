package zio

import (
	"context"
	"errors"
	"io"
)

func CtxReader(ctx context.Context, r io.Reader) io.Reader {
	return &ctxReader{r: r, ctx: ctx}
}

type readResult struct {
	n   int
	err error
}

type ctxReader struct {
	ctx    context.Context
	r      io.Reader
	cancel func()
	sig    chan readResult
}

func (r *ctxReader) Read(p []byte) (n int, err error) {
	if r.sig == nil {
		r.sig = make(chan readResult, 1)
	}
	go r.goRead(p, r.sig)
	select {
	case <-r.ctx.Done():
		if r.cancel != nil {
			r.cancel()
		}
		if err := r.ctx.Err(); err != nil {
			return 0, err
		}
		return 0, errors.New("context canceled without error")
	case res := <-r.sig:
		if r.cancel != nil {
			r.cancel()
		}
		return res.n, res.err
	}
}

func (r *ctxReader) goRead(p []byte, ch chan readResult) {
	n, err := r.r.Read(p)
	ch <- readResult{n, err}
}
