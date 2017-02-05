package store

import (
	"time"
)

type unixTime int64

func (u unixTime) Human() string {
	return time.Unix(int64(u), 0).Format("02/01/2006 15:04")
}

func UnixTime(i int64) unixTime {
	return unixTime(i)
}

func UnixNow() unixTime {
	return unixTime(time.Now().Unix())
}
