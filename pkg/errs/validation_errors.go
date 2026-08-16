package errs

import (
	"errors"
)

var (
	ErrInvalidValue                  = errors.New("недопустимое значение")
	ErrRequiredField                 = errors.New(`это обязательное поле`)
	ErrInvalidData                   = errors.New("переданы некорректные данные")
	ErrValidationFailed              = errors.New("невалидный JSON или не переданы необходимые поля")
	ErrPathParametrized              = errors.New("параметры пути содержат недопустимые символы")
	ErrInvalidAmount                 = errors.New("указана некорректная сумма")
	ErrInvalidID                     = errors.New("указан некорректный идентификатор")
	ErrInvalidTitle                  = errors.New("заголовок содержит недопустимые символы или имеет неверную длину")
	ErrInvalidToken                  = errors.New("токен авторизации недействителен или истек")
	ErrUserIDIsEmpty                 = errors.New("идентификатор пользователя не может быть пустым")
	ErrInvalidDescription            = errors.New("описание содержит недопустимые символы или имеет неверную длину")
	ErrInvalidEmail                  = errors.New("указан некорректный формат email адреса")
	ErrInvalidFullName               = errors.New("имя должно содержать от 2 до 100 символов и состоять только из букв и пробелов")
	ErrFullNameIsEmpty               = errors.New("полное имя не может быть пустым")
	ErrInvalidRegion                 = errors.New("название региона должно содержать от 2 до 100 символов")
	ErrRegionIsEmpty                 = errors.New("регион не может быть пустым")
	ErrInvalidGender                 = errors.New("указано недопустимое значение пола (допустимые: MALE, FEMALE, OTHER)")
	ErrInvalidAge                    = errors.New("возраст должен быть в диапазоне от 18 до 120 лет")
	ErrInvalidMaritalStatus          = errors.New("указано недопустимое значение семейного положения (допустимые: SINGLE, MARRIED, DIVORCED, WIDOWED)")
	ErrInvalidRole                   = errors.New("указана недопустимая роль (допустимые: admin, viewer, experimenter, approver)")
	ErrRoleMustBeExperimenterOrAdmin = errors.New("роль должна быть experimenter либо admin")
	ErrInvalidToFromDate             = errors.New("указаны недопустимые значения для фильтрации по датам")

	ErrInvalidPag            = errors.New("указаны неверные данные для пагинации")
	ErrInvalidDateFormat     = errors.New("дата должна быть в формате RFC3339")
	ErrInvalidDateRange      = errors.New("from должен быть меньше to")
	ErrDateRangeTooLarge     = errors.New("разница между from и to не должна превышать 90 дней")
	ErrInvalidStatus         = errors.New("статус может быть только APPROVED или DECLINED")
	ErrValidationFailedQuery = errors.New("невалидные фильтры(query)")
)
