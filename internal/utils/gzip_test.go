package utils

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompressWriterCloseWithoutBodyDoesNotWriteGZIP(t *testing.T) {
	response := httptest.NewRecorder()
	writer := NewGZIPCompressWriter(response)

	writer.WriteHeader(http.StatusNoContent)
	require.NoError(t, writer.Close())

	res := response.Result()
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	assert.Empty(t, res.Header.Get("Content-Encoding"))
	assert.Empty(t, body)
}

func TestCompressWriterWritesGZIPBody(t *testing.T) {
	response := httptest.NewRecorder()
	writer := NewGZIPCompressWriter(response)

	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write([]byte("hello"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	res := response.Result()
	defer res.Body.Close()

	gzipReader, err := gzip.NewReader(res.Body)
	require.NoError(t, err)
	defer gzipReader.Close()

	body, err := io.ReadAll(gzipReader)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
	assert.Equal(t, "hello", string(body))
}
