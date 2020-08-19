package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToLineSetCurrent(t *testing.T) {
	var demoLineFromHandle = LineFromHandle{map[string]string{"FOOTBALL": "123"}}
	var incorrectLineFromHandle = LineFromHandle{map[string]string{"FOOTBALL": "123", "BASEBALL": "456"}}
	var expectedLineSetCurrent = LineSetCurrent{LineName: "football", LineCurrentValue: 123}

	res, err := demoLineFromHandle.ConvertToLineSetCurrent()
	assert.Equal(t, expectedLineSetCurrent, *res)
	assert.Nil(t, err)

	res, err = incorrectLineFromHandle.ConvertToLineSetCurrent()
	assert.Nil(t, res)
	assert.Error(t, err)
}
