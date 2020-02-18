package util

import (
	"strconv"
	"time"
)

func GetDateToday() int {
	return GetDate(time.Now())
}

func GetDate(t time.Time) int {
	n, err := strconv.Atoi(t.Format("20060102"))
	if err != nil {
		panic(err) // 不废话，直接panic
	}
	return n
}
