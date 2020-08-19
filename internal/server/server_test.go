package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BiteLikeASnake/kiddy_line_processor/internal/model"
	"github.com/stretchr/testify/assert"
)

////тест aliveHandler при успешной синхронизации с lines_provider
func Test_aliveHandlerReady(t *testing.T) {
	model.ResponsesFromLinesCounter = 3
	req, err := http.NewRequest("GET", "/ready", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(aliveHandler)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

//тест aliveHandler без синхронизации с lines_provider
func Test_aliveHandlerNotReady(t *testing.T) {
	model.ResponsesFromLinesCounter = 0
	req, err := http.NewRequest("GET", "/ready", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(aliveHandler)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
