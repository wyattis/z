package zhttp

import (
	"bytes"
	"io"
	"io/fs"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/wyattis/z/zio"
)

func TestFileInterfaces(t *testing.T) {
	var _ http.File = (*HttpFile)(nil)
	var _ fs.File = (*HttpFile)(nil)
	var _ io.ReadSeekCloser = (*HttpFile)(nil)
}

func TestHttpFileReadAll(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dir, _ := filepath.Abs("../test/assets")
		http.ServeFile(w, r, filepath.Join(dir, r.URL.Path))
	}))
	httpFile := NewHttpFile(ts.URL + "/test.mp4")
	defer httpFile.Close()
	file, err := os.Open("../test/assets/test.mp4")
	defer file.Close()
	if err != nil {
		t.Error(err)
	}
	if err = zio.ReadersMatch(httpFile, file, 8000); err != nil {
		t.Error(err)
	}
}

func TestHttpFileSeek(t *testing.T) {
	data := make([]byte, 5*32*1000)
	if _, err := io.ReadFull(rand.New(rand.NewSource(1000)), data); err != nil {
		t.Error(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, "none", time.Now(), bytes.NewReader(data))
	}))

	httpFile := NewHttpFile(ts.URL)
	defer httpFile.Close()
	if _, err := httpFile.Seek(3000, io.SeekEnd); err != nil {
		t.Error(err)
	}
	if err := zio.ReadersMatch(httpFile, bytes.NewReader(data[len(data)-3000:]), 200); err != nil {
		t.Error("seek 3000 left", err)
	}

}
