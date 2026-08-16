package date_time

import "time"

func GetCurrentTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}
