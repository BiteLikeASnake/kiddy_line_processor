package server

import (
	"fmt"
	"net/http"

	"github.com/BiteLikeASnake/kiddy_line_processor/internal/model"
	"github.com/gorilla/mux"
)

//Connector - содержит роутер и адрес вызываемого сервиса
type Connector struct {
	router *mux.Router
	addr   string
}

//Конструктор *Connector
func New(addr string) *Connector {
	c := &Connector{}
	c.router = mux.NewRouter()
	c.addr = addr
	return c
}

func (c *Connector) executeHandlers() {
	c.router.HandleFunc("/ready", aliveHandler).Methods("GET")
}

//Start запуск http сервера
func (c *Connector) Start() {
	c.executeHandlers()
	go http.ListenAndServe(c.addr, c.router)
}

//aliveHandler проверяет что пришли ответы от сервиса по всем линиям
func aliveHandler(w http.ResponseWriter, r *http.Request) {
	if model.ResponsesFromLinesCounter < model.LinesAmount {
		http.Error(w, fmt.Sprintf("Got responses from %d of %d lines", model.ResponsesFromLinesCounter, model.LinesAmount), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Hello from kiddy line processor")
}
