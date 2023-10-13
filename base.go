package goutils

import "time"

type NowUnixTimeI interface {
	Unix() int64
}

type NowUnixTime struct {
}

func (*NowUnixTime) Unix() int64 {
	return time.Now().Unix()
}

type MockUnixTime struct {
	value int64
}

func (t *MockUnixTime) Unix() int64 {
	return t.value
}

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func BoolToInt(v bool) int {
	if v {
		return 1
	} else {
		return 0
	}
}

func IntToBool(v int) bool {
	return v != 0
}
