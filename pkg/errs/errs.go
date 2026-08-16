package errs

type ErrorCode string

const (
	ErrorCodeBadRequest            ErrorCode = "BAD_REQUEST"
	ErrorCodeValidationFailed      ErrorCode = "VALIDATION_FAILED"
	ErrorCodeUnauthorized          ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden             ErrorCode = "FORBIDDEN"
	ErrorCodeNotFound              ErrorCode = "NOT_FOUND"
	ErrorCodeUserNotFound          ErrorCode = "USER_NOT_FOUND"
	ErrorCodeAlreadyExists         ErrorCode = "ALREADY_EXISTS"
	ErrorCodeEmailAlreadyExists    ErrorCode = "EMAIL_ALREADY_EXISTS"
	ErrorCodeUserInactive          ErrorCode = "USER_INACTIVE"
	ErrorCodeDslParseError         ErrorCode = "DSL_PARSE_ERROR"
	ErrorCodeDslInvalidField       ErrorCode = "DSL_INVALID_FIELD"
	ErrorCodeDslInvalidOperator ErrorCode = "DSL_INVALID_OPERATOR"
	ErrorCodeConflict           ErrorCode = "CONFLICT"
	ErrorCodeInternalServerError   ErrorCode = "INTERNAL_SERVER_ERROR"
)
