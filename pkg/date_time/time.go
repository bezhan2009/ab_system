package date_time

import (
	"time"
)

func ParseRFC3339(str string) (time.Time, error) {
	if str == "" {
		return time.Time{}, nil
	}

	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
