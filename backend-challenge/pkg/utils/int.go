package utils

import (
	"fmt"
	"strconv"
)

func GetInt64(n interface{}) (int64, error) {
	switch i := n.(type) {
	case int64:
		return i, nil
	case int32:
		return int64(i), nil
	case int:
		return int64(i), nil
	case uint64:
		return int64(i), nil
	case uint32:
		return int64(i), nil
	case uint:
		return int64(i), nil
	case string:
		return strconv.ParseInt(i, 10, 64)
	case float64:
		return int64(i), nil
	case float32:
		return int64(i), nil
	}

	return 0, fmt.Errorf("cant convert n=%v to int64", n)
}
