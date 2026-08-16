package errs

import "errors"

var (
	ErrMissingCache = errors.New("кэш не найден")
)
