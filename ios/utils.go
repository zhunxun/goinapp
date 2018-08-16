package ios

import "time"

// convertToTime convert unix timestamp in milliseconds to Go time.Time
func convertToTime(timeMS int64) time.Time {
	return time.Unix(0, timeMS*int64(time.Millisecond))
}
