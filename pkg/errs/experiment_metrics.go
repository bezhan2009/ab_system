package errs

import "errors"

var (
	ErrMultiplePrimaryMetrics = errors.New("может быть только одна основная метрика")
	ErrThresholdRequired      = errors.New("для guardrail требуется порог")
	ErrOperatorRequired       = errors.New("для guardrail требуется оператор")
	ErrInvalidOperator        = errors.New("недопустимый оператор (допустимы: >, >=, <, <=)")
	ErrWindowMinRequired      = errors.New("для guardrail требуется окно наблюдения")
	ErrActionRequired         = errors.New("для guardrail требуется действие")
	ErrInvalidAction          = errors.New("недопустимое действие (допустимы: pause, rollback)")
)
