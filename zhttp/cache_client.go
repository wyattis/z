package zhttp

import (
	"bytes"
	"encoding/gob"
	"io"
	"net/http"

	"github.com/wyattis/z/zio"
)

// func NewCacheClient(cache zcache.Cache) *CacheClient {
// 	return &CacheClient{}
// }

// This aims to be a replacement for the default http.Client with customizable
// file caching for requests.
// type CacheClient struct {
// 	http.Client
// 	RequestFilter  func(*http.Request) bool
// 	ResponseFilter func(*http.Response) bool
// 	HashFunc       func(*http.Request) string
// 	cache          zcache.Cache
// }

// func (c *CacheClient) hash(req *http.Request) string {
// 	if c.HashFunc != nil {
// 		return c.HashFunc(req)
// 	}
// 	return zpath.FileEscape(req.URL.String())
// }

// func (c *CacheClient) shouldUseCachedRequest(req *http.Request) bool {
// 	if c.RequestFilter != nil {
// 		return c.RequestFilter(req)
// 	}
// 	return req.Method == http.MethodGet
// }

// func (c *CacheClient) shouldCacheResponse(res *http.Response) bool {
// 	if c.ResponseFilter != nil {
// 		return c.ResponseFilter(res)
// 	}
// 	return res.Request.Method == http.MethodGet && res.StatusCode == 200
// }

// func (c *CacheClient) respondWithCached(key string, req *http.Request) (res *http.Response, err error) {
// 	f, err := c.cache.Open(key)
// 	if err != nil {
// 		return
// 	}
// 	defer f.Close()
// 	cRes := cacheResponse{}
// 	if err = cRes.Decode(f); err != nil {
// 		return
// 	}
// 	res = cRes.AsResponse()
// 	return
// }

// func (c *CacheClient) Do(req *http.Request) (res *http.Response, err error) {
// 	// just return without caching this result
// 	if !c.shouldUseCachedRequest(req) {
// 		return c.Client.Do(req)
// 	}
// 	key := c.hash(req)
// 	if c.cache.Has(key) {
// 		return c.respondWithCached(key, req)
// 	}
// 	res, err = c.Client.Do(req)
// 	if err != nil {
// 		return
// 	}
// 	// just return the response if we don't want to cache this response
// 	if !c.shouldCacheResponse(res) {
// 		return
// 	}
// 	w, err := c.cache.Create(key)
// 	if err != nil {
// 		return
// 	}
// 	cRes := cacheResponse{}
// 	cRes.SetResponse(res)
// 	res.Body = zio.CombineReaderCloser{
// 		Reader: io.TeeReader(res.Body, bytes.NewBuffer(cRes.BodyBuf)),
// 		Closer: zio.MultiCloser(res.Body, cRes.EncodeOnClose(w)),
// 	}
// 	return
// }

// a serializable version of http.Request including the body
type cacheResponse struct {
	Status        string
	StatusCode    int
	Header        http.Header
	BodyBuf       []byte
	Proto         string
	ContentLength int64

	res *http.Response
}

func (r *cacheResponse) Close() error {
	if r.res != nil {
		return zio.CloseAll(r.res.Body)
	}
	r.BodyBuf = nil
	return nil
}

func (r *cacheResponse) SetResponse(res *http.Response) {
	r.Header = res.Header.Clone()
	r.Status = res.Status
	r.StatusCode = res.StatusCode
	r.ContentLength = res.ContentLength
	r.Proto = res.Proto
	r.res = res
}

func (r cacheResponse) AsResponse() *http.Response {
	return &http.Response{
		Status:        r.Status,
		StatusCode:    r.StatusCode,
		Header:        r.Header.Clone(),
		Body:          &zio.ReaderToReadCloser{Reader: bytes.NewBuffer(r.BodyBuf)},
		Proto:         r.Proto,
		ContentLength: r.ContentLength,
	}
}

func (r cacheResponse) Encode(writer io.Writer) (err error) {
	dec := gob.NewEncoder(writer)
	return dec.Encode(r)
}

func (r *cacheResponse) Decode(reader io.Reader) (err error) {
	enc := gob.NewDecoder(reader)
	return enc.Decode(&r)
}

func (r *cacheResponse) EncodeOnClose(writer io.WriteCloser) io.Closer {
	return &cacheResponseCloser{res: r, writer: writer}
}

type cacheResponseCloser struct {
	res    *cacheResponse
	writer io.WriteCloser
}

func (c *cacheResponseCloser) Close() error {
	defer c.writer.Close()
	return c.res.Encode(c.writer)
}
