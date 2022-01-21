package zhttp

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrNoRangeSupport       = errors.New("remote server doesn't support seeking")
	ErrUnsupportedRangeType = errors.New("unsupported range type")
	ErrNegativeSeek         = errors.New("cannot seek below 0")
	ErrExceededLength       = errors.New("exceeded content size")
	ErrNoWrite              = errors.New("unable to write to http server")
)

type Range struct {
	Start int
	End   int
	Total int
}

func parseRange(header http.Header) (r Range, err error) {
	if header.Get("Content-Range") != "" {
		parts := strings.Split(header.Get("Content-Range"), " ")
		units, rParts := parts[0], parts[1]
		if units != "bytes" {
			err = ErrUnsupportedRangeType
			return
		}
		rangeParts := strings.Split(rParts, "/")
		parts, totalStr := strings.Split(rangeParts[0], "-"), rangeParts[1]
		var start, end, total int64
		start, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return
		}
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return
		}
		total, err = strconv.ParseInt(totalStr, 10, 64)
		if err != nil {
			return
		}
		r.Start = int(start)
		r.End = int(end)
		r.Total = int(total)
	}
	// TODO: Maybe we should just remove this
	if r.Total == 0 && header.Get("Content-Length") != "" {
		length, err := strconv.ParseUint(header.Get("Content-Length"), 10, 64)
		if err == nil {
			r.Total = int(length)
		}
	}
	return
}
