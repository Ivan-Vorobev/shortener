package handler

import (
	"Ivan-Vorobev/shortener/internal/config"
	logger "Ivan-Vorobev/shortener/internal/logger"
	"Ivan-Vorobev/shortener/internal/model"
	"Ivan-Vorobev/shortener/internal/utils"
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/buffer"
)

func newTestConfig(t *testing.T) *config.Configuration {
	t.Helper()

	conf := config.NewDefaultConfig()
	conf.FileStoragePath = filepath.Join(t.TempDir(), conf.FileStoragePath)

	return conf
}

func newTestRouter(t *testing.T, conf *config.Configuration) http.Handler {
	t.Helper()

	log, err := logger.NewLogger()
	require.NoError(t, err)

	shutdown := utils.NewShutdown()
	t.Cleanup(func() {
		require.NoError(t, shutdown.Close())
	})

	router, err := NewRouter(conf, log, shutdown)
	require.NoError(t, err)

	return router
}

func TestCreateLink(t *testing.T) {
	link := "https://yandex.ru"
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(link))
	request.Header.Set("Content-Type", "text/plain")

	response := httptest.NewRecorder()

	conf := newTestConfig(t)
	router := newTestRouter(t, conf)
	router.ServeHTTP(response, request)

	res := response.Result()
	// проверяем код ответа
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// получаем и проверяем тело запроса
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)

	assert.NoError(t, err)

	responseURL := fmt.Sprintf("%s/", conf.BaseURL)
	require.NoError(t, err)
	assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
	assert.True(t, strings.HasPrefix(string(resBody), responseURL))
	assert.True(t, len(resBody) > len(responseURL))
}

func TestAPICreateLink(t *testing.T) {
	link := `{"url":"https://yandex.ru"}`
	request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(link))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	conf := newTestConfig(t)
	router := newTestRouter(t, conf)
	router.ServeHTTP(response, request)

	res := response.Result()
	// проверяем код ответа
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// получаем и проверяем тело запроса
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)

	assert.NoError(t, err)

	outURL := model.OutURL{}
	err = json.Unmarshal(resBody, &outURL)
	assert.NoError(t, err)

	responseURL := fmt.Sprintf("%s/", conf.BaseURL)
	require.NoError(t, err)

	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	assert.True(t, strings.HasPrefix(string(outURL.Result), responseURL))
	assert.True(t, len(outURL.Result) > len(responseURL))
}

func TestAPICreateGZIPLinkRequest(t *testing.T) {
	var body buffer.Buffer
	link := `{"url":"https://yandex.ru"}`

	compressedBody := gzip.NewWriter(&body)
	compressedBody.Write([]byte(link))
	compressedBody.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(body.String()))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Accept-Encoding", "gzip")

	response := httptest.NewRecorder()

	conf := newTestConfig(t)
	router := newTestRouter(t, conf)
	router.ServeHTTP(response, request)

	res := response.Result()
	// проверяем код ответа
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// Проверяем что ответ пришел в сжатом виде
	assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))

	// получаем и проверяем тело запроса
	defer res.Body.Close()
	gzipReader, err := gzip.NewReader(res.Body)
	assert.NoError(t, err)
	defer gzipReader.Close()

	resBody, err := io.ReadAll(gzipReader)

	assert.NoError(t, err)

	outURL := model.OutURL{}
	err = json.Unmarshal(resBody, &outURL)
	assert.NoError(t, err)

	responseURL := fmt.Sprintf("%s/", conf.BaseURL)
	require.NoError(t, err)

	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	assert.True(t, strings.HasPrefix(string(outURL.Result), responseURL))
	assert.True(t, len(outURL.Result) > len(responseURL))
}

func TestGetLink(t *testing.T) {
	// Сохраняем полный url
	link := "https://yandex.ru"
	requestPost := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(link))
	requestPost.Header.Set("Content-Type", "text/plain; charset=utf-8")

	responsePost := httptest.NewRecorder()

	conf := newTestConfig(t)
	router := newTestRouter(t, conf)
	router.ServeHTTP(responsePost, requestPost)

	res := responsePost.Result()
	// проверяем код ответа
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// Получаем короткий url и проверяем тело запроса
	defer res.Body.Close()
	shortLink, err := io.ReadAll(res.Body)

	assert.NoError(t, err)

	shortLinkURL, err := url.Parse(string(shortLink))
	assert.NoError(t, err)
	// Проверяем что у нас короткий url без домена и queryString
	assert.True(t, strings.HasSuffix(string(shortLink), shortLinkURL.Path))

	// собираем второй запрос с короткой ссылкой, чтобы проверить что нам вернется наша полная ссылка
	requestGet := httptest.NewRequest(http.MethodGet, shortLinkURL.Path, nil)
	requestGet.Header.Set("Content-Type", "text/plain; charset=utf-8")

	responseGet := httptest.NewRecorder()

	router.ServeHTTP(responseGet, requestGet)

	resGet := responseGet.Result()

	// получаем и проверяем тело запроса
	defer resGet.Body.Close()
	fullLink, err := io.ReadAll(resGet.Body)
	assert.NoError(t, err)

	assert.Equal(t, "text/plain", resGet.Header.Get("Content-Type"))
	assert.Equal(t, http.StatusTemporaryRedirect, resGet.StatusCode)
	assert.Equal(t, link, resGet.Header.Get("Location"))
	assert.Empty(t, fullLink)
}

func TestAPIStorageCreateLink(t *testing.T) {
	link := "https://yandex.ru"
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/shorten",
		strings.NewReader(fmt.Sprintf(`{"url":"%s"}`, link)))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	conf := newTestConfig(t)
	router := newTestRouter(t, conf)
	router.ServeHTTP(response, request)

	res := response.Result()
	// проверяем код ответа
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// получаем и проверяем тело запроса
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)

	assert.NoError(t, err)

	outURL := model.OutURL{}
	err = json.Unmarshal(resBody, &outURL)
	assert.NoError(t, err)

	file, err := os.OpenFile(conf.FileStoragePath, os.O_RDONLY, os.ModePerm)
	assert.NoError(t, err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	assert.True(t, scanner.Scan())

	var linkStorage model.LinkStorageRow
	line := scanner.Text()

	if err := json.Unmarshal([]byte(line), &linkStorage); err != nil {
		assert.NoError(t, err)
	}

	originalLink, err := linkStorage.OriginalURL.Get()
	assert.NoError(t, err)

	assert.Equal(t, "1", linkStorage.UUID)
	assert.Equal(t, link, originalLink)
	assert.Equal(t, outURL.Result, fmt.Sprintf("%s%s", conf.BaseURL, linkStorage.ShortURL.String()))
}

func TestAPIRestoreCreateLink(t *testing.T) {
	conf := newTestConfig(t)
	shortLink := "Vu2LiXCO"
	originalLink, err := model.NewLink("https://yandex.ru")

	assert.NoError(t, err)

	file, err := os.OpenFile(conf.FileStoragePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	assert.NoError(t, err)
	writer := bufio.NewWriter(file)

	storageData := model.LinkStorageRow{
		UUID:        "1",
		ShortURL:    model.NewShortLink(shortLink),
		OriginalURL: originalLink,
	}

	storageDataJSON, err := json.Marshal(storageData)
	storageDataJSON = append(storageDataJSON, byte('\n'))
	t.Log(string(storageDataJSON))
	assert.NoError(t, err)
	bl, err := writer.WriteString(string(storageDataJSON))
	assert.NotEqual(t, 0, bl)
	assert.NoError(t, err)
	require.NoError(t, writer.Flush())
	assert.NoError(t, file.Close())

	requestGet := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", shortLink), nil)
	requestGet.Header.Set("Content-Type", "text/plain; charset=utf-8")

	responseGet := httptest.NewRecorder()

	router := newTestRouter(t, conf)
	router.ServeHTTP(responseGet, requestGet)

	resGet := responseGet.Result()

	// получаем и проверяем тело запроса
	defer resGet.Body.Close()
	fullLink, err := io.ReadAll(resGet.Body)
	assert.NoError(t, err)

	redirectURL, err := originalLink.Get()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, resGet.StatusCode)
	assert.Equal(t, redirectURL, resGet.Header.Get("Location"))
	assert.Empty(t, fullLink)
}
