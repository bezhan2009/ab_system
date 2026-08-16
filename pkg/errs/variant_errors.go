package errs

import "errors"

var (
	ErrNoVariants                 = errors.New("эксперимент должен содержать хотя бы один вариант")
	ErrInvalidVariantWeight       = errors.New("вес варианта должен быть положительным числом")
	ErrVariantNameRequired        = errors.New("название варианта обязательно")
	ErrInvalidControlVariantCount = errors.New("эксперимент должен содержать ровно один контрольный вариант")
	ErrVariantWeightSumMismatch   = errors.New("сумма весов вариантов должна равняться проценту трафика эксперимента")
	ErrVariantNotFound            = errors.New("вариант не найден")
)
