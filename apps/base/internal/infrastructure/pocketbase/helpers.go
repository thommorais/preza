package pocketbase

import "time"

func ms(start time.Time) int64 {
	return time.Since(start).Milliseconds()
}
