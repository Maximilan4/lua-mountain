package attr

import (
	"fmt"
	"time"
)

func GetTyped[RValue any](m map[string]any, key string) (val RValue, ok bool) {
	if _, ok = m[key]; ok {
		val, ok = m[key].(RValue)
	}
	return
}

func GetDuration(m map[string]any, key string) (val time.Duration, err error) {
	if sDur, ok := GetTyped[string](m, key); ok {
		return time.ParseDuration(sDur)
	}

	if iDur, ok := GetTyped[int64](m, key); ok {
		return time.Duration(iDur), nil
	}

	return 0, fmt.Errorf("unsupported duration value by key %s", key)
}
