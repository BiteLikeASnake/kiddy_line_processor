package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/BiteLikeASnake/kiddy_line_processor/internal/grpcserver"
	"github.com/BiteLikeASnake/kiddy_line_processor/internal/model"
	"github.com/BiteLikeASnake/kiddy_line_processor/internal/server"
	"github.com/BiteLikeASnake/kiddy_line_processor/internal/storage"
	"github.com/labstack/gommon/log"
)

//var provider_address = "http://localhost:8000/api/v1/lines"
//var storage_address = "user=postgres password=example dbname=lines_storage sslmode=disable port=5432 host=localhost"

var client = &http.Client{
	Timeout: time.Second * 10,
}

var config *model.Config = &model.Config{}

const (
	timeFootbalSec = 1
)

func main() {
	log.Print("Started")
	//Устанавливаем значения переменных окружения
	err := config.GetConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	//Открываем http соединение
	server := server.New(config.HttpPort)
	server.Start()

	//Подключаем базу данных
	model.Storage, err = storage.New(config.StorageConn)
	if err != nil {
		log.Fatal(err.Error())
	}
	//Запускаем функции обращения к ручкам сервиса lines_provider
	go autoCall(config.SInterval, "soccer")
	go autoCall(config.BInterval, "baseball")
	go autoCall(config.FInterval, "football")
	//Открываем gRPC соединение
	err = runGrpcServer(config.GrpcPort)
	if err != nil {
		log.Fatal(err.Error())
	}
}

//autoCall обращается к lines_provider, сохраняет полученное значение в хранилище, ждет необходимое время.
func autoCall(delay int, suffix string) {
	firstTimeCalled := true

	finalAddress := fmt.Sprintf(config.ProviderAddress + "/" + suffix)
	request, err := http.NewRequest("GET", finalAddress, nil)
	if err != nil {
		log.Print(err.Error())
	}
	linesFromHandle := &model.LineFromHandle{}
	for {
		err := updateStorageWithLine(linesFromHandle, request)
		if err != nil {
			log.Error("autocall", err)
		} else if firstTimeCalled {
			model.ResponsesFromLinesCounter++
			firstTimeCalled = false
		}
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

//updateStorageWithLine сохраняет значение в хранилище
func updateStorageWithLine(lineFromHandle *model.LineFromHandle, request *http.Request) error {
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("updateStorageWithLine: %v", err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("updateStorageWithLine: %v", err)
	}
	err = json.Unmarshal(res, lineFromHandle)
	if err != nil {
		return fmt.Errorf("updateStorageWithLine: %v", err)
	}
	setCurrent, err := lineFromHandle.ConvertToLineSetCurrent()
	if err != nil {
		return fmt.Errorf("updateStorageWithLine: %s", err.Error())
	}
	err = model.Storage.UpdateLineCurrentVal(*setCurrent)
	if err != nil {
		return fmt.Errorf("updateStorageWithLine: %s", err.Error())
	}
	return nil
}

//runGrpcServer запускает gRPC сервер
func runGrpcServer(port string) error {
	for {
		if model.ResponsesFromLinesCounter == model.LinesAmount {
			err := grpcserver.StartServer(port)
			return err
		}
		time.Sleep(time.Duration(500) * time.Millisecond)
	}
}
