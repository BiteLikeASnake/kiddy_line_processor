package model

import (
	"fmt"
	"strconv"
	"strings"
)

//LineFromHandle ...
type LineFromHandle struct {
	Lines map[string]string `json:"lines"`
}

//LineSetCurrent ...
type LineSetCurrent struct {
	LineName         string  `gorm:"primary_key;column:line_name"`
	LineCurrentValue float64 `gorm:"column:line_current_value"`
}

//Lines ...
type Lines struct {
	LineName         string  `gorm:"primary_key;column:line_name"`
	LineCurrentValue float64 `gorm:"column:line_current_value"`
	LineLatestValue  float64 `gorm:"column:line_latest_value"`
}

//ConvertToLineSetCurrent ...
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

//TODO config, 2ступенчатое получение данных
//Envs ...
type Envs struct {
	HttpPort        string `long:"http" env:"HTTP_PORT" description:"Http port" default:":8081"`
	GrpcPort        string `long:"grpc" env:"GRPC_PORT" description:"grpc port" default:":9000"`
	StorageConn     string `long:"stconn" env:"STORAGE" description:"Connection string to storage database" default:"user=postgres password=example dbname=lines_storage sslmode=disable port=5432 host=localhost"`
	ProviderAddress string `long:"provider" env:"PROVIDER_ADDRESS" description:"lines provider http address" default:"http://localhost:8000/api/v1/lines"`
	FInterval       int    `long:"fi" env:"FOOTBALL_INTERVAL" description:"time interval in seconds for football handle" default:"10"`
	SInterval       int    `long:"si" env:"SOCCER_INTERVAL" description:"time interval in seconds for soccer handle" default:"10"`
	BInterval       int    `long:"bi" env:"BASEBALL_INTERVAL" description:"time interval in seconds for baseball handle" default:"10"`
}
