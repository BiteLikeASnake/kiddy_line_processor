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
	"github.com/jessevdk/go-flags"
	"github.com/labstack/gommon/log"
)

//var provider_address = "http://localhost:8000/api/v1/lines"
//var storage_address = "user=postgres password=example dbname=lines_storage sslmode=disable port=5432 host=localhost"

var client = &http.Client{
	Timeout: time.Second * 10,
}

var envs model.Envs
var parser = flags.NewParser(&envs, flags.Default)
var err error

const (
	timeFootbalSec = 1
)

func main() {
	//Устанавливаем значения переменных окружения
	if _, err := parser.Parse(); err != nil {
		log.Fatal(err.Error())
	}
	//Открываем http соединение
	server := server.New(envs.HttpPort)
	server.Start()

	//Подключаем базу данных
	model.Storage, err = storage.New(envs.StorageConn)
	if err != nil {
		log.Fatal(err.Error)
	}
	//Запускаем функции обращения к ручкам сервиса lines_provider
	go autoCall(envs.SInterval, "soccer")
	go autoCall(envs.BInterval, "baseball")
	go autoCall(envs.FInterval, "football")
	//Открываем gRPC соединение
	runGrpcServer(envs.GrpcPort)
}

//autoCall обращается к lines_provider, сохраняет полученное значение в хранилище, ждет необходимое время.
func autoCall(delay int, suffix string) {
	firstTimeCalled := true

	finalAddress := fmt.Sprintf(envs.ProviderAddress + "/" + suffix)
	request, err := http.NewRequest("GET", finalAddress, nil)
	if err != nil {
		fmt.Println(err.Error())
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
