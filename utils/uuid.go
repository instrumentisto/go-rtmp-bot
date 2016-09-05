package utils

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"time"
)

// Returns unique string.
func GetUUID() string {
	// generate 32 bits timestamp
	unix32bits := uint32(time.Now().UTC().Unix())
	buff := make([]byte, 12)
	numRead, err := rand.Read(buff)
	if numRead != len(buff) || err != nil {
		panic(err)
	}
	return fmt.Sprintf(
		"%x-%x-%x-%x-%x-%x",
		unix32bits,
		buff[0:2],
		buff[2:4],
		buff[4:6],
		buff[6:8],
		buff[8:])
}

func Num64(n interface{}) int64 {
	s := fmt.Sprintf("%d", n)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	} else {
		return i
	}
}
