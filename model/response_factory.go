package model

import (
	"bytes"
	"encoding/json"
)

// Creates response for status request.

var status_descriptions = map[uint8]string{
	0: "Test ready. Click start for starts this.",
	1: "Test running. Click stop for stops this.",
	2: "Test error.Ups! Something wrong...",
}

// Returns new status result.
//
// param status uint8   Status of test.
func GetResponse(status uint8) string {
	result := NewResult(status, status_descriptions[status])
	jsn, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(jsn).String()
}
