package util

import (
	"qnhd/pkg/logging"
	"strconv"
)

func AsUint(a string) uint64 {
	b, err := strconv.ParseUint(a, 10, 64)
	if err != nil {
		logging.Error(err.Error())
		panic(err)
	}
	return b
}

func AsInt(a string) int {
	b := AsUint(a)
	return int(b)
}
