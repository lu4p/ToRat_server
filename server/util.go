package server

import (
	"time"
)

func getTimeSt() string {
	now := time.Now()
	return now.Format("2006-01-02_15:04:05")
}
