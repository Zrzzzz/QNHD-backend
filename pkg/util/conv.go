package util

import (
	"fmt"
	"qnhd/pkg/logging"
	"strconv"
)

// string转uint64
func AsUint(a string) uint64 {
	b, err := strconv.ParseUint(a, 10, 64)
	if err != nil {
		logging.Error(err.Error())
		panic(err)
	}
	return b
}

// string转int
func AsInt(a string) int {
	b := AsUint(a)
	return int(b)
}

// uint转string
func AsStrU(a uint64) string {
	return fmt.Sprintf("%d", a)
}

// int转string
func AsStr(a int) string {
	return fmt.Sprintf("%d", a)
}
