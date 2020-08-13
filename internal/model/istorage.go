package model

//IStorage интерфейс хранилища данных
type IStorage interface {
	UpdateLineCurrentVal(LineSetCurrent) error
	ReturnLineDelta(line string) (float64, error)
	ReturnLineCurrentVal(line string) (float64, error)
}
