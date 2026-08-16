package errs

import "errors"

var (
	ErrUserInactive = errors.New("учетная запись пользователя деактивирована")
	ErrUserNotFound = errors.New("пользователь не найден")
)
