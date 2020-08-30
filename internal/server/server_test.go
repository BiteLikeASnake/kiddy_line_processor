package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/call-me-snake/kiddy_line_processor/internal/model"
	"github.com/stretchr/testify/assert"
)

//TestReadyHandlerReady - тест readyHandler при успешной синхронизации с lines_provider
func TestReadyHandlerReady(t *testing.T) {
	model.ResponsesFromLinesCounter = 3
	req, err := http.NewRequest("GET", "/ready", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(readyHandler)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

//TestReadyHandlerNotReady - тест readyHandler без синхронизации с lines_provider
func TestReadyHandlerNotReady(t *testing.T) {
	model.ResponsesFromLinesCounter = 0
	req, err := http.NewRequest("GET", "/ready", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(readyHandler)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
