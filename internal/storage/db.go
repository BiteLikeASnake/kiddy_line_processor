package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/call-me-snake/kiddy_line_processor/internal/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//sleepDuration - время пинга функции checkConnection в секундах
const sleepDuration = 5

//storage ...
type storage struct {
	database *gorm.DB
	Address  string
}

//New возвращает объект интерфейса IStorage (storage)
func New(adress string) (model.IStorage, error) {
	var err error
	db := &storage{}
	db.Address = adress
	db.database, err = gorm.Open("postgres", adress)
	if err != nil {
		return nil, fmt.Errorf("db.New: %v", err)
	}

	err = db.ping()
	if err != nil {
		return nil, fmt.Errorf("db.New: %s", err.Error())
	}
	db.checkConnection()

	return db, nil
}

//ping (internal)
func (db *storage) ping() error {
	//db.database.LogMode(true)
	result := struct {
		Result int
	}{}

	err := db.database.Raw("select 1+1 as result").Scan(&result).Error
	if err != nil {
		return fmt.Errorf("db.ping: %v", err)
	}
	if result.Result != 2 {
		return fmt.Errorf("db.ping: incorrect result!=2 (%d)", result.Result)
	}
	return nil
}

//checkConnection (internal)
func (db *storage) checkConnection() {
	go func() {
		for {
			err := db.ping()
			if err != nil {

				log.Printf("db.checkConnection: no connection: %s", err.Error())
				tempDb, err := gorm.Open("postgres", db.Address)

				if err != nil {
					log.Printf("db.checkConnection: could not establish connection: %v", err)
				} else {
					db.database = tempDb
				}
			}
			time.Sleep(sleepDuration * time.Second)
		}
	}()
}

//UpdateLineCurrentVal - обновляет поле LineCurrentValue. Используется при получении данных из сервиса Lines_provider
func (db *storage) UpdateLineCurrentVal(setVal model.LineSetCurrent) error {

	//db.database.LogMode(true)
	update := db.database.Model(model.Lines{}).Where("line_name = ?", setVal.LineName).Update(setVal)

	if update.Error != nil {
		return fmt.Errorf("SetLineCurrentVal: %v", update.Error)
	}
	return nil
}

//ReturnLineDelta - используется grpc сервером. Возвращает разность между LineCurrentValue и line_latest_value. Также обновляет line_latest_value до LineCurrentValue в таблице.
func (db *storage) ReturnLineDelta(line string) (float64, error) {
	queryLine := model.Lines{}
	query := db.database.Where("line_name = ?", line).First(&queryLine)
	if query.Error != nil {
		return 0, fmt.Errorf("ReturnLineDelta: %v", query.Error)
	}
	if query.RowsAffected == 0 {
		return 0, fmt.Errorf("ReturnLineDelta: не найдено строки с line_name = %s", line)
	}
	delta := queryLine.LineCurrentValue - queryLine.LineLatestValue
	update := db.database.Model(&model.Lines{}).Where("line_name = ?", line).Update("line_latest_value", queryLine.LineCurrentValue)
	if update.Error != nil {
		return 0, fmt.Errorf("ReturnLineDelta: %v", update.Error)
	}
	if update.RowsAffected == 0 {
		return 0, fmt.Errorf("ReturnLineDelta: не найдено строки с line_name = %s", line)
	}
	return delta, nil
}

//ReturnLineCurrentVal - используется grpc сервером. Возвращает LineCurrentValue для строки с line_name = line. Также обновляет line_latest_value до LineCurrentValue в таблице.
func (db *storage) ReturnLineCurrentVal(line string) (float64, error) {
	queryLine := model.Lines{}
	query := db.database.Where("line_name = ?", line).First(&queryLine)
	if query.Error != nil {
		return 0, fmt.Errorf("ReturnLineCurrentVal: %v", query.Error)
	}
	if query.RowsAffected == 0 {
		return 0, fmt.Errorf("ReturnLineCurrentVal: не найдено строки с line_name = %s", line)
	}
	update := db.database.Model(&model.Lines{}).Where("line_name = ?", line).Update("line_latest_value", queryLine.LineCurrentValue)
	if update.Error != nil {
		return 0, fmt.Errorf("UpdateLineLatestVal: %v", update.Error)
	}
	if update.RowsAffected == 0 {
		return 0, fmt.Errorf("UpdateLineLatestVal: не найдено строки с line_name = %s", line)
	}
	return queryLine.LineCurrentValue, nil
}
