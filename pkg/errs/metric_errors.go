package errs

import "errors"

var (
	ErrInvalidMetricType   = errors.New("недопустимый тип метрики")
	ErrInvalidMetricConfig = errors.New("некорректная конфигурация метрики для выбранного типа")
)
