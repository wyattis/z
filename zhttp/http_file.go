package zhttp

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strconv"
	"time"
)

type HttpFileInfo struct {
	source HttpFile
}

func (f *HttpFileInfo) Name() string {
	return f.source.Url
}

func (f *HttpFileInfo) Size() int64 {
	return f.source.ContentSize
}

func (f *HttpFileInfo) Mode() fs.FileMode {
	return fs.ModeDevice
}

func (f *HttpFileInfo) ModTime() time.Time {
	return f.source.ModTime
}

func (f *HttpFileInfo) IsDir() bool {
	return false
}

func (f *HttpFileInfo) Sys() interface{} {
	return nil
}

type HttpFile struct {
	Url       string
	Requests  uint64
	ReadBytes uint64
	isOpen    bool
	offset    int64
	client    *http.Client
	response  *http.Response

	ContentType  string
	ContentSize  int64
	RangeSupport bool
	ModTime      time.Time

	io.ReadSeekCloser
	io.ReaderAt
}

func NewHttpFile(url string) *HttpFile {
	return &HttpFile{
		Url: url,
	}
}

func (f *HttpFile) SetClient(client *http.Client) {
	f.client = client
}

func (f *HttpFile) Name() string {
	return f.Url
}

func (f *HttpFile) ReadAt(p []byte, offset int64) (n int, err error) {
	// fmt.Println("readAt", offset, len(p), f.ContentSize)
	if err = f.Open(); err != nil {
		return
	}
	if f.RangeSupport {
		return f.readRange(p, offset)
	} else {
		return 0, ErrNoRangeSupport
	}
}

func (f *HttpFile) Read(p []byte) (n int, err error) {
	if err = f.Open(); err != nil {
		return
	}
	// fmt.Println("reading", len(p), f.offset, f.ContentSize)
	if f.RangeSupport {
		n, err = f.readRange(p, f.offset)
		f.offset += int64(n)
		return
	}
	return f.response.Body.Read(p)
}

func (f *HttpFile) Stat() (fs.FileInfo, error) {
	if err := f.Open(); err != nil {
		return nil, err
	}
	info := &HttpFileInfo{
		source: *f,
	}
	return info, nil
}

func (f *HttpFile) Readdir(count int) ([]fs.FileInfo, error) {
	// TODO: Implement directory listing over HTTP
	res := []fs.FileInfo{}
	return res, nil
}

func (f *HttpFile) Readdirnames(count int) ([]string, error) {
	// TODO: Implement directory listing over HTTP
	return []string{}, nil
}

func (f *HttpFile) readRange(p []byte, offset int64) (n int, err error) {
	req, err := http.NewRequest("GET", f.Url, nil)
	if err != nil {
		return
	}
	end := offset + int64(len(p))
	if end > f.ContentSize {
		end = f.ContentSize
	}
	rangeHeader := fmt.Sprintf("bytes=%d-%d", offset, end)
	req.Header.Set("Range", rangeHeader)
	if err != nil {
		return
	}
	res, err := (*f.client).Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusPartialContent {
		return 0, ErrNoRangeSupport
	}
	contentRange, err := parseRange(res.Header)
	if err != nil {
		return n, err
	}
	rangeSize := contentRange.End - contentRange.Start
	n = len(p)
	if rangeSize < n {
		n = rangeSize + 1
	}

	n, err = io.ReadFull(res.Body, p)
	if err != nil {
		if errors.Is(err, io.ErrUnexpectedEOF) {
			err = nil
		} else {
			return
		}
	}
	// Finish reading this request so the connection can be reused
	if len(p) > n {
		io.Copy(io.Discard, res.Body)
	}
	offset += int64(n)
	if offset >= f.ContentSize {
		err = io.EOF
	}
	f.Requests += 1
	f.ReadBytes += uint64(n)
	return
}

func (f *HttpFile) Open() (err error) {
	if f.isOpen {
		return
	}

	if f.client == nil {
		f.client = http.DefaultClient
	}

	req, err := http.NewRequest(http.MethodGet, f.Url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Range", "bytes=0-")
	// fmt.Println("init req", req.URL, req.Method, req.Header)
	res, err := (*f.client).Do(req)
	if err != nil {
		return
	}
	// fmt.Println("res", res.StatusCode)
	if res.StatusCode == http.StatusPartialContent {
		f.RangeSupport = true
		defer res.Body.Close()
	} else {
		f.response = res
	}
	contentLengthStr := res.Header.Get("Content-Length")
	if contentLengthStr != "" {
		total, err := strconv.ParseInt(contentLengthStr, 10, 64)
		if err != nil {
			return err
		}
		f.ContentSize = total
	}
	f.isOpen = true
	f.Requests += 1
	return
}

func (f *HttpFile) Seek(offset int64, whence int) (int64, error) {
	// fmt.Println("seeking", offset, whence)
	if err := f.Open(); err != nil {
		return 0, err
	}

	if !f.RangeSupport {
		return 0, ErrNoRangeSupport
	}

	// TODO: handle all the other whences
	if whence == io.SeekEnd {
		offset = f.ContentSize - offset
	} else if whence == io.SeekCurrent {
		offset = int64(f.offset) + offset
	}

	if offset < 0 {
		return 0, ErrNegativeSeek
	} else if offset > f.ContentSize {
		return 0, ErrExceededLength
	}

	f.offset = offset
	return f.offset, nil
}

func (f *HttpFile) Close() error {
	if f.response != nil {
		return f.response.Body.Close()
	}
	return nil
}

// All write methods unavailable
func (f *HttpFile) Write(p []byte) (int, error)                 { return 0, ErrNoWrite }
func (f *HttpFile) WriteAt(p []byte, offset int64) (int, error) { return 0, ErrNoWrite }
func (f *HttpFile) Sync() error                                 { return ErrNoWrite }
func (f *HttpFile) Truncate(size int64) error                   { return ErrNoWrite }
func (f *HttpFile) WriteString(str string) (int, error)         { return 0, ErrNoWrite }
