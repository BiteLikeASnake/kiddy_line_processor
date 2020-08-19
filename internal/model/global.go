package model

//Storage - переменная хранилища
var Storage IStorage

//ResponsesFromLinesCounter - счетчик линий lines_provider, по которым поступают ответы.
var ResponsesFromLinesCounter int = 0

//LinesAmount - количество линий lines_provider, к которым поступают запросы.
const LinesAmount = 3

/*
//LinesMap ...
var LinesMap = map[string]bool{
	"football": true,
	"soccer":   true,
	"baseball": true,
}
*/
