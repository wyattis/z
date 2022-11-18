package ztext

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/gogs/chardet"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
)

var boms = map[string][]byte{
	"UTF-8":    {0xef, 0xbb, 0xbf},
	"UTF-16LE": {0xff, 0xfe},
	"UTF-16BE": {0xfe, 0xff},
	"UTF-32BE": {0x00, 0x00, 0xfe, 0xff},
	"UTF-32LE": {0xff, 0xfe, 0x00, 0x00},
}

func NewUTF8ReaderCharset(source io.Reader, charset string) *UTF8Reader {
	return &UTF8Reader{
		source:  source,
		charset: charset,
		reader:  bufio.NewReader(source),
	}
}

func NewUTF8Reader(source io.Reader) *UTF8Reader {
	return &UTF8Reader{
		source: source,
		reader: bufio.NewReader(source),
	}
}

type UTF8Reader struct {
	source         io.Reader
	reader         *bufio.Reader
	charset        string
	decodingReader io.Reader
}

func (r *UTF8Reader) deleteBom(head []byte, charset string) {
	bom, ok := boms[charset]
	if !ok {
		return
	}
	if bytes.Equal(bom, head[:len(bom)]) {
		r.reader.Discard(len(bom))
	}
}

func (r *UTF8Reader) init() (err error) {
	if r.decodingReader != nil {
		return
	}

	head, err := r.reader.Peek(2048)
	// Try to autodetect the encoding
	if r.charset == "" {
		if err == io.EOF {
			err = nil
		} else if err != nil {
			return err
		}
		encResult, err := chardet.NewTextDetector().DetectBest(head)
		if err != nil {
			return err
		}
		r.charset = encResult.Charset
	}

	r.deleteBom(head, r.charset)

	if r.charset == "UTF-8" {
		r.decodingReader = r.reader
		return
	}

	e, err := ianaindex.MIME.Encoding(r.charset)
	if err != nil {
		return err
	} else if e == nil {
		return fmt.Errorf("unsupported charset: %s", r.charset)
	}
	r.decodingReader = transform.NewReader(r.reader, e.NewDecoder())
	return
}

func (r *UTF8Reader) Read(d []byte) (n int, err error) {
	if err = r.init(); err != nil {
		return
	}
	n, err = r.decodingReader.Read(d)
	return
}

func (r *UTF8Reader) Close() (err error) {
	if c, ok := r.source.(io.Closer); ok {
		return c.Close()
	}
	return
}
