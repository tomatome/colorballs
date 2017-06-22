// util.go
package util

import (
	"encoding/json"
	"log"
	"strconv"
)

type Msg struct {
	OpCode int
	Data   interface{}
}

func Encoder(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		log.Println("[ERROR]: failed to encoding<", err, ">")
		return ""
	}

	return string(data)
}

func Decoder(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		log.Println("[ERROR]: Failed to decode<", err, ">")
	}
}
func Atoi(a string) int {
	b, err := strconv.Atoi(a)
	if err != nil {
		log.Println("[ERROR]: Atoi ", err)
		return 0
	}
	return b
}
