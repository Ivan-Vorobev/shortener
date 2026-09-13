package utils

import (
	"compress/gzip"
	"io"
	"net/http"
	"slices"
	"strings"

	"go.uber.org/zap/buffer"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func NewGZIPCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 || statusCode >= 400 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressForHeaderWriter struct {
	compressWriter
	allowedContentTypes []string
	canGZIP             bool
}

func NewGZIPCompressWriterForHeader(w http.ResponseWriter, allowedContentTypes []string) *compressForHeaderWriter {
	cw := NewGZIPCompressWriter(w)
	return &compressForHeaderWriter{
		compressWriter:      *cw,
		allowedContentTypes: allowedContentTypes,
		canGZIP:             false,
	}
}

func (c *compressForHeaderWriter) WriteHeader(statusCode int) {
	if statusCode < 300 || statusCode >= 400 {
		contentType := strings.ToLower(c.w.Header().Get("Content-Type"))
		var contentTypes []string
		if strings.Contains(contentType, ";") {
			contentTypes = strings.SplitN(contentType, ";", 2)
		} else {
			contentTypes = []string{contentType}
		}
		c.canGZIP = slices.Contains(c.allowedContentTypes, contentTypes[0])

		if c.canGZIP {
			c.w.Header().Set("Content-Encoding", "gzip")
		}
	}
	c.w.WriteHeader(statusCode)
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func (c *compressForHeaderWriter) Write(p []byte) (int, error) {
	if c.canGZIP {
		return c.zw.Write(p)
	}

	var b buffer.Buffer
	c.zw.Reset(&b)
	return c.w.Write(p)
}

func NewGZIPCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
