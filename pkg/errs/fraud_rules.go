package errs

import "errors"

var (
	ErrInvalidFraudName          = errors.New("диапазон допустимых символов название правил 3..120")
	ErrInvalidFraudPriority      = errors.New("приоритет правила должно быть >= 1")
	ErrFraudNameUniquenessFailed = errors.New("правило с таким именем уже существует")
)
