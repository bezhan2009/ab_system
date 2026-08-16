package errs

import "errors"

var (
	ErrDuplicateEvent = errors.New("событие с таким event_id уже существует")
)
