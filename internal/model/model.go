package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
)

//LineFromHandle - структура для хранения ответа ручек lines_provider
type LineFromHandle struct {
	Lines map[string]string `json:"lines"`
}

//LineSetCurrent - структура для хранения текущего значения линии в формате gorm
type LineSetCurrent struct {
	LineName         string  `gorm:"primary_key;column:line_name"`
	LineCurrentValue float64 `gorm:"column:line_current_value"`
}

//Lines - структура для хранения полей линии в формате gorm
type Lines struct {
	LineName         string  `gorm:"primary_key;column:line_name"`
	LineCurrentValue float64 `gorm:"column:line_current_value"`
	LineLatestValue  float64 `gorm:"column:line_latest_value"`
}

//ConvertToLineSetCurrent - функция преобразования структуры LineFromHandle в LineSetCurrent
func (val LineFromHandle) ConvertToLineSetCurrent() (*LineSetCurrent, error) {
	if len(val.Lines) != 1 {
		return nil, fmt.Errorf("ConvertToLineSetCurrent: Входной параметр %+v имеет неверный формат.\n", val)
	}
	res := &LineSetCurrent{}
	var err error
	for key, value := range val.Lines {
		res.LineName = strings.ToLower(key)
		res.LineCurrentValue, err = strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, fmt.Errorf("ConvertToLineSetCurrent: Ошибка преобразования строки %s в числовой тип: %s\n", value, err.Error())
		}
	}
	return res, nil
}

//envs - приватная структура, получает переменные окружения
type envs struct {
	HttpPort        string `long:"http" env:"HTTP_PORT" description:"Http port" default:":8081"`
	GrpcPort        string `long:"grpc" env:"GRPC_PORT" description:"grpc port" default:":9000"`
	StorageConn     string `long:"stconn" env:"STORAGE" description:"Connection string to storage database" default:"user=postgres password=example dbname=lines_storage sslmode=disable port=5432 host=localhost"`
	ProviderAddress string `long:"provider" env:"PROVIDER_ADDRESS" description:"lines provider http address" default:"http://localhost:8000/api/v1/lines"`
	FInterval       string `long:"fi" env:"FOOTBALL_INTERVAL" description:"time interval in seconds for football handle" default:"10"`
	SInterval       string `long:"si" env:"SOCCER_INTERVAL" description:"time interval in seconds for soccer handle" default:"10"`
	BInterval       string `long:"bi" env:"BASEBALL_INTERVAL" description:"time interval in seconds for baseball handle" default:"10"`
}

//Config - публичная структура, хранит проверенные переменные окружения
type Config struct {
	HttpPort        string
	GrpcPort        string
	StorageConn     string
	ProviderAddress string
	FInterval       int
	SInterval       int
	BInterval       int
}

//GetConfig - получает переменные окружения с помощью приватной структуры envs и проверяет их
func (c *Config) GetConfig() error {
	e := envs{}
	parser := flags.NewParser(&e, flags.Default)
	if _, err := parser.Parse(); err != nil {
		return fmt.Errorf("GetConfig: %v", err)
	}
	c.HttpPort = e.HttpPort
	c.GrpcPort = e.GrpcPort
	c.StorageConn = e.StorageConn
	c.ProviderAddress = e.ProviderAddress

	val, err := strconv.Atoi(e.FInterval)
	if err != nil {
		return fmt.Errorf("GetConfig: %v", err)
	}
	if val <= 0 {
		return fmt.Errorf("GetConfig: Получено значение FInterval <=0")
	}
	c.FInterval = val

	val, err = strconv.Atoi(e.SInterval)
	if err != nil {
		return fmt.Errorf("GetConfig: %v", err)
	}
	if val <= 0 {
		return fmt.Errorf("GetConfig: Получено значение SInterval <=0")
	}
	c.SInterval = val

	val, err = strconv.Atoi(e.BInterval)
	if err != nil {
		return fmt.Errorf("GetConfig: %v", err)
	}
	if val <= 0 {
		return fmt.Errorf("GetConfig: Получено значение BInterval <=0")
	}
	c.BInterval = val
	return nil
}
