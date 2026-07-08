package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateLink(t *testing.T) {
	link := "https://yandex.ru"
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(link))
	request.Header.Set("Content-Type", "text/plain")

	response := httptest.NewRecorder()

	handler := CreateMainHandler()
	handler(response, request)

	res := response.Result()
	// проверяем код ответа
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// получаем и проверяем тело запроса
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)

	responseUrl := fmt.Sprintf("http://%s/", request.Host)
	require.NoError(t, err)
	assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
	assert.True(t, strings.HasPrefix(string(resBody), responseUrl))
	assert.True(t, len(resBody) > len(responseUrl))
}

func TestGetLink(t *testing.T) {
	// Сохраняем полный url
	link := "https://yandex.ru"
	requestPost := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(link))
	requestPost.Header.Set("Content-Type", "text/plain; charset=utf-8")

	responsePost := httptest.NewRecorder()

	handler := CreateMainHandler()
	handler(responsePost, requestPost)

	res := responsePost.Result()
	// проверяем код ответа
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// Получаем короткий url и проверяем тело запроса
	defer res.Body.Close()
	shortLink, err := io.ReadAll(res.Body)

	assert.NoError(t, err)

	shortLinkUrl, err := url.Parse(string(shortLink))
	assert.NoError(t, err)
	// Проверяем что у нас короткий url без домена и queryString
	assert.True(t, strings.HasSuffix(string(shortLink), shortLinkUrl.Path))

	// собираем второй запрос с короткой ссылкой, чтобы проверить что нам вернется наша полная ссылка
	requestGet := httptest.NewRequest(http.MethodGet, shortLinkUrl.Path, nil)
	requestGet.Header.Set("Content-Type", "text/plain; charset=utf-8")

	responseGet := httptest.NewRecorder()

	handler(responseGet, requestGet)

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
