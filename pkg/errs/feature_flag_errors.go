package errs

import "errors"

var (
	ErrFeatureFlagAlreadyExists = errors.New("флаг с таким ключом уже существует")
	ErrFeatureFlagNotFound      = errors.New("флаг не найден")
)
