package model

var Storage IStorage

var ResponsesFromLinesCounter int = 0

const LinesAmount = 3

//LinesMap ...
var LinesMap = map[string]bool{
	"football": true,
	"soccer":   true,
	"baseball": true,
}
