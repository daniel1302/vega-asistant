package utils

import (
	"fmt"
	"strconv"
)

func MustUint64(val string) uint64 {
	if val == "" {
		val = "0"
	}

	uVal, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("failed to parse uint value from string: %s", err.Error()))
	}

	return uVal
}
