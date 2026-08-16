package errs

import "errors"

var (
	ErrCantGenerateTrace  = errors.New("не удалось сгенерировать trace ID для запроса")
	ErrSomethingWentWrong = errors.New("произошла непредвиденная ошибка")
)
