package utils

import (
	"compress/gzip"
	"io"
	"net/http"
	"slices"
	"strings"
)

type compressWriter struct {
	w           http.ResponseWriter
	zw          *gzip.Writer
	wroteHeader bool
	canGZIP     bool
	compressed  bool
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
	if !c.wroteHeader {
		c.WriteHeader(http.StatusOK)
	}

	if !c.canGZIP {
		return c.w.Write(p)
	}

	return c.writeGZIP(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if c.wroteHeader {
		return
	}

	c.wroteHeader = true
	c.canGZIP = canCompressStatus(statusCode)

	if c.canGZIP {
		c.w.Header().Set("Content-Encoding", "gzip")
	}

	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	if !c.compressed {
		return nil
	}

	return c.zw.Close()
}

func (c *compressWriter) writeGZIP(p []byte) (int, error) {
	n, err := c.zw.Write(p)
	if n > 0 || err == nil {
		c.compressed = true
	}

	return n, err
}

type compressForHeaderWriter struct {
	compressWriter
	allowedContentTypes []string
}

func NewGZIPCompressWriterForHeader(w http.ResponseWriter, allowedContentTypes []string) *compressForHeaderWriter {
	cw := NewGZIPCompressWriter(w)
	return &compressForHeaderWriter{
		compressWriter:      *cw,
		allowedContentTypes: allowedContentTypes,
	}
}

func (c *compressForHeaderWriter) WriteHeader(statusCode int) {
	if c.wroteHeader {
		return
	}

	c.wroteHeader = true
	c.canGZIP = false

	if canCompressStatus(statusCode) {
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
	if !c.wroteHeader {
		c.WriteHeader(http.StatusOK)
	}

	if c.canGZIP {
		return c.writeGZIP(p)
	}

	return c.w.Write(p)
}

func (c *compressForHeaderWriter) Close() error {
	return c.compressWriter.Close()
}

func canCompressStatus(statusCode int) bool {
	if statusCode >= 100 && statusCode < 200 {
		return false
	}

	if statusCode == http.StatusNoContent || statusCode == http.StatusNotModified {
		return false
	}

	return statusCode < 300 || statusCode >= 400
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
