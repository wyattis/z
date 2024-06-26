package zio

import (
	"context"
	"io"
	"time"
)

func TimeoutReader(reader io.Reader, timeout time.Duration) io.Reader {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &ctxReader{r: reader, ctx: ctx, cancel: cancel}
}
