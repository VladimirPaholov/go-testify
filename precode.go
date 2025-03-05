package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req := httptest.NewRequest("GET", "/cafe?count=5&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Корректный статус 200 OK")

	// Проверяем, что вернулись ВСЕ кафе
	expectedResponse := strings.Join(cafeList["moscow"], ",")
	assert.Equal(t, expectedResponse, rec.Body.String(), "Корректный ответ - все доступные кафе")

	// Проверяем, что количество кафе в ответе соответствует totalCount
	cafes := strings.Split(rec.Body.String(), ",")
	assert.Len(t, cafes, totalCount, "Корректный ответ - все кафе")

}

// Отправляем корректный запрос с count=2 и city=moscow
func TestMainHandlerValidRequest(t *testing.T) {

	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", nil)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Корректный статус 200 OK")

	assert.NotEmpty(t, rec.Body.String(), "Тело ответа не должно быть пустым")
}

// Город, который передаётся в параметре city, не поддерживается.
func TestMainHandlerWithInvalidCity(t *testing.T) {

	req := httptest.NewRequest("GET", "/cafe?count=2&city=unknown", nil)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(rec, req)

	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code, "Корректный статус 400 Bad Request")

	assert.Equal(t, "wrong city value", rec.Body.String(), "Вывод - 'wrong city value'")
}
