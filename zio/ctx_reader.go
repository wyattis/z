package zio

import (
	"context"
	"io"
)

func CtxReader(ctx context.Context, r io.Reader) io.Reader {
	return &ctxReader{r: r, ctx: ctx}
}

type ctxReader struct {
	ctx context.Context
	r   io.Reader
}

func (r *ctxReader) Read(p []byte) (n int, err error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.r.Read(p)
}
