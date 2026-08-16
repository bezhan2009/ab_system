package errs

import "errors"

var (
	ErrInvalidCredentials    = errors.New("неверные учетные данные")
	ErrPasswordIsEmpty       = errors.New("пароль не может быть пустым")
	ErrUsernameIsEmpty       = errors.New("имя пользователя не может быть пустым")
	ErrInvalidPhoneNumber    = errors.New("указан неверный формат номера телефона")
	ErrEmailIsEmpty          = errors.New("email не может быть пустым")
	ErrEmailUniquenessFailed = errors.New("пользователь с таким email уже зарегистрирован в системе")
	ErrUnauthorized          = errors.New("требуется авторизация для доступа к данному ресурсу")
)
