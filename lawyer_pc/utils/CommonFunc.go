package utils

import (
	"strconv"
	"time"
)

func FormatTime(timestamp string) string {
	temp, _ := strconv.Atoi(timestamp)
	t := time.Unix(int64(temp), 0)
	dateStr := t.Format("2006-01-02 15:04:05")
	return dateStr
}
