package zhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/wyattis/z/zslice/zstrings"
	"github.com/wyattis/z/zstring"
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

// Decode the body into a struct using the Content-Type header
func DecodeRequestBody(dest interface{}, r *http.Request) error {
	fullType := r.Header.Get("Content-Type")
	cType, _, _ := zstring.Cut(fullType, ";")
	switch cType {
	case "text/json":
		fallthrough
	case "application/json":
		dec := json.NewDecoder(r.Body)
		return dec.Decode(dest)
	default:
		return mapstructure.Decode(r.Form, dest)
	}
}

// Encode the body as json
func Json(w http.ResponseWriter, value interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	if err := enc.Encode(value); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// Encode StatusBadRequest with optional msg
func BadRequest(w http.ResponseWriter, msg string) {
	if msg == "" {
		msg = http.StatusText(http.StatusBadRequest)
	}
	http.Error(w, msg, http.StatusBadRequest)
}

// Encode StatusInternalServerError with optional msg
func InternalServerError(w http.ResponseWriter, msg string) {
	if msg == "" {
		msg = http.StatusText(http.StatusInternalServerError)
	}
	http.Error(w, msg, http.StatusInternalServerError)
}

type MultipartHandler = func(w http.ResponseWriter, r *http.Request, file multipart.File, header *multipart.FileHeader) error

// Process multipart upload
func MultipartUpload(maxSize int64, allowedTypes []string, handler MultipartHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > maxSize {
			BadRequest(w, "exceeded max size")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxSize)
		if err := r.ParseMultipartForm(maxSize); err != nil {
			BadRequest(w, "exceeded max size")
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			BadRequest(w, err.Error())
			return
		}

		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			InternalServerError(w, err.Error())
			return
		}

		// TODO: any way to make this more reliable?
		// filetype := http.DetectContentType(buff)
		filetype := fileHeader.Header.Get("Content-Type")
		if !zstrings.Contains(allowedTypes, filetype) {
			fmt.Println(allowedTypes, filetype)
			BadRequest(w, "the provided file format is not allowed")
			return
		}

		if _, err = file.Seek(0, io.SeekStart); err != nil {
			InternalServerError(w, err.Error())
			return
		}

		if err = handler(w, r, file, fileHeader); err != nil {
			InternalServerError(w, err.Error())
			return
		}
	}
}
