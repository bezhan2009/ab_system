package errs

import "errors"

var (
	ErrRecordNotFound   = errors.New("запись не найдена в базе данных")
	ErrPermissionDenied = errors.New("недостаточно прав для выполнения данной операции")
	ErrIdIsInvalid      = errors.New("неверный UUID")
	ErrAlreadyExists    = errors.New("ресурс уже существует")
	ErrDeleteFailed     = errors.New("не удалось удалить запись из базы данных")
)
