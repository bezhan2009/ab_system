package response

import (
	"ab_system/internal/http/middlewares/observability"
	"ab_system/pkg/date_time"
	"ab_system/pkg/errs"
	"ab_system/pkg/logger"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var errorMappings = map[error]errs.ErrorMapping{
	// 400 Bad Request
	errs.ErrInvalidData:                      {http.StatusBadRequest, errs.ErrorCodeBadRequest, nil},
	errs.ErrValidationFailed:                 {http.StatusBadRequest, errs.ErrorCodeBadRequest, map[string]interface{}{"hint": "Проверьте запятые/кавычки и необходимые поля"}},
	errs.ErrPathParametrized:                 {http.StatusBadRequest, errs.ErrorCodeBadRequest, nil},
	errs.ErrInvalidVariantWeight:             {http.StatusBadRequest, errs.ErrorCodeBadRequest, map[string]interface{}{"field": "variants[].weight", "hint": "Вес варианта должен быть положительным числом"}},
	errs.ErrVariantNameRequired:              {http.StatusBadRequest, errs.ErrorCodeBadRequest, map[string]interface{}{"field": "variants[].name", "hint": "Название варианта обязательно"}},
	errs.ErrNoVariants:                       {http.StatusBadRequest, errs.ErrorCodeBadRequest, map[string]interface{}{"field": "variants", "hint": "Эксперимент должен содержать хотя бы один вариант"}},
	errs.ErrInvalidControlVariantCount:       {http.StatusBadRequest, errs.ErrorCodeBadRequest, map[string]interface{}{"field": "variants", "hint": "Должен быть ровно один контрольный вариант"}},
	errs.ErrVariantWeightSumMismatch:         {http.StatusBadRequest, errs.ErrorCodeBadRequest, map[string]interface{}{"field": "variants", "hint": "Сумма весов должна равняться traffic_percent"}},
	errs.ErrExperimentStatusIsNotOnReview:    {http.StatusBadRequest, errs.ErrorCodeBadRequest, map[string]interface{}{"hint": "отправьте эксперимент на ревю"}},
	errs.ErrExperimentAlreadyHasBeenOnReview: {http.StatusBadRequest, errs.ErrorCodeBadRequest, nil},
	errs.ErrInvalidStatusTransition:          {http.StatusBadRequest, errs.ErrorCodeBadRequest, nil},
	errs.ErrExperimentNotApproved:            {http.StatusBadRequest, errs.ErrorCodeBadRequest, nil},
	errs.ErrConclusionRequired:               {http.StatusBadRequest, errs.ErrorCodeBadRequest, nil},
	errs.ErrRoleMustBeExperimenterOrAdmin:    {http.StatusBadRequest, errs.ErrorCodeBadRequest, nil},
	errs.ErrInvalidMetricType:                {http.StatusBadRequest, errs.ErrorCodeBadRequest, nil},
	errs.ErrInvalidMetricConfig:              {http.StatusBadRequest, errs.ErrorCodeBadRequest, nil},
	errs.ErrThresholdRequired:                {http.StatusBadRequest, errs.ErrorCodeBadRequest, nil},
	errs.ErrOperatorRequired:                 {http.StatusBadRequest, errs.ErrorCodeBadRequest, nil},
	errs.ErrInvalidOperator:                  {http.StatusBadRequest, errs.ErrorCodeBadRequest, nil},
	errs.ErrWindowMinRequired:                {http.StatusBadRequest, errs.ErrorCodeBadRequest, nil},
	errs.ErrActionRequired:                   {http.StatusBadRequest, errs.ErrorCodeBadRequest, nil},
	errs.ErrInvalidAction:                    {http.StatusBadRequest, errs.ErrorCodeBadRequest, nil},
	errs.ErrInvalidWinnerVariant:             {http.StatusBadRequest, errs.ErrorCodeBadRequest, nil},

	// 401 Unauthorized
	errs.ErrInvalidToken:       {http.StatusUnauthorized, errs.ErrorCodeUnauthorized, nil},
	errs.ErrUnauthorized:       {http.StatusUnauthorized, errs.ErrorCodeUnauthorized, nil},
	errs.ErrInvalidCredentials: {http.StatusUnauthorized, errs.ErrorCodeUnauthorized, nil},

	// 403 Forbidden
	errs.ErrPermissionDenied:            {http.StatusForbidden, errs.ErrorCodeForbidden, nil},
	errs.ErrCannotArchiveNotCompleted:   {http.StatusForbidden, errs.ErrorCodeForbidden, map[string]interface{}{"hint": "Можно архивировать только завершенные эксперименты"}},
	errs.ErrForbiddenToBeApprover:       {http.StatusForbidden, errs.ErrorCodeForbidden, nil},
	errs.ErrApproverAlreadyVoted:        {http.StatusForbidden, errs.ErrorCodeForbidden, nil},
	errs.ErrCannotEditRunningExperiment: {http.StatusForbidden, errs.ErrorCodeForbidden, map[string]interface{}{"hint": "Эксперимент находится в режиме \"Заморозка\""}},

	// 404 Not Found
	errs.ErrRecordNotFound:      {http.StatusNotFound, errs.ErrorCodeNotFound, nil},
	errs.ErrUserNotFound:        {http.StatusNotFound, errs.ErrorCodeNotFound, nil},
	errs.ErrTransactionNotFound: {http.StatusNotFound, errs.ErrorCodeNotFound, nil},
	errs.ErrVariantNotFound:     {http.StatusNotFound, errs.ErrorCodeNotFound, nil},
	errs.ErrFeatureFlagNotFound: {http.StatusNotFound, errs.ErrorCodeNotFound, nil},

	// 409 Conflict
	errs.ErrEmailUniquenessFailed:     {http.StatusConflict, errs.ErrorCodeEmailAlreadyExists, map[string]interface{}{"field": "email"}},
	errs.ErrFeatureFlagAlreadyExists:  {http.StatusConflict, errs.ErrorCodeAlreadyExists, map[string]interface{}{"field": "key"}},
	errs.ErrFraudNameUniquenessFailed: {http.StatusConflict, errs.ErrorCodeAlreadyExists, map[string]interface{}{"field": "title"}},
	errs.ErrDuplicateEntry:            {http.StatusConflict, errs.ErrorCodeAlreadyExists, nil},
	errs.ErrAlreadyExists:             {http.StatusConflict, errs.ErrorCodeAlreadyExists, map[string]interface{}{"hint": "Ресурс с таким названием/параметрами уже существует"}},
	errs.ErrActiveExperimentExists:    {http.StatusConflict, errs.ErrorCodeConflict, map[string]interface{}{"field": "flag_key", "hint": "На этот флаг уже есть активный эксперимент"}},
	errs.ErrDuplicateEvent:            {http.StatusConflict, errs.ErrorCodeAlreadyExists, nil},

	// 422 Unprocessable entity
	errs.ErrValidationFailedQuery: {http.StatusUnprocessableEntity, errs.ErrorCodeBadRequest, map[string]interface{}{"hint": "Проверьте ваш URL query"}},

	// 423 Locked
	errs.ErrUserInactive: {http.StatusLocked, errs.ErrorCodeUserInactive, nil},
}

var validationErrorMappings = map[error]errs.FieldError{
	errs.ErrPasswordIsEmpty:      {Field: "password"},
	errs.ErrEmailIsEmpty:         {Field: "email"},
	errs.ErrInvalidEmail:         {Field: "email"},
	errs.ErrFullNameIsEmpty:      {Field: "fullName"},
	errs.ErrInvalidFullName:      {Field: "fullName"},
	errs.ErrRegionIsEmpty:        {Field: "region"},
	errs.ErrInvalidRegion:        {Field: "region"},
	errs.ErrInvalidAge:           {Field: "age"},
	errs.ErrInvalidGender:        {Field: "gender"},
	errs.ErrInvalidMaritalStatus: {Field: "maritalStatus"},
	errs.ErrInvalidRole:          {Field: "role"},
	errs.ErrUsernameIsEmpty:      {Field: "username"},
	errs.ErrInvalidPhoneNumber:   {Field: "phone"},
	errs.ErrInvalidAmount:        {Field: "amount"},
	errs.ErrInvalidID:            {Field: "id"},
	errs.ErrInvalidTitle:         {Field: "title"},
	errs.ErrInvalidDescription:   {Field: "description"},
	errs.ErrUserIDIsEmpty:        {Field: "userId"},
	errs.ErrInvalidField:         {Field: "general"},
	errs.ErrInvalidToFromDate:    {Field: "to, from"},
	errs.ErrInvalidDateRange:     {Field: "to, from"},
	errs.ErrDateRangeTooLarge:    {Field: "to, from"},
	errs.ErrInvalidDateFormat:    {Field: "to, from"},
	errs.ErrInvalidStatus:        {Field: "status"},
	errs.ErrInvalidFraudName:     {Field: "name"},
	errs.ErrInvalidFraudPriority: {Field: "priority"},
	errs.ErrIdIsInvalid:          {Field: "id"},
	errs.ErrInvalidPag:           {Field: "size, page"},
}

func getPathFromContext(c *gin.Context) string {
	return c.Request.URL.Path
}

func isValidationError(err error) bool {
	_, ok1 := err.(*errs.ValidationError)
	_, ok2 := err.(*errs.MultiValidationError)
	_, ok3 := err.(*errs.ValidationErrors)
	return ok1 || ok2 || ok3
}

func mapErrorToResponse(err error, traceId, timestamp, path string) (int, interface{}) {
	if isValidationError(err) {
		return handleValidationError(err, traceId, timestamp, path)
	}

	if nf, ok := err.(*errs.NotFoundWithDetailsError); ok {
		return http.StatusNotFound, errs.ApiErrorResponse{
			Code:      errs.ErrorCodeNotFound,
			Message:   nf.Message,
			TraceId:   traceId,
			Timestamp: timestamp,
			Path:      path,
			Details:   nf.Details,
		}
	}

	for mappedErr, mapping := range errorMappings {
		if errors.Is(err, mappedErr) {
			return mapping.HttpStatus, errs.ApiErrorResponse{
				Code:      mapping.ErrorCode,
				Message:   err.Error(),
				TraceId:   traceId,
				Timestamp: timestamp,
				Path:      path,
				Details:   mapping.Details,
			}
		}
	}

	for validationErr, fieldError := range validationErrorMappings {
		if errors.Is(err, validationErr) {
			fieldError.Issue = validationErr.Error()
			return http.StatusUnprocessableEntity, errs.ValidationErrorResponse{
				Code:        errs.ErrorCodeValidationFailed,
				Message:     "Некоторые поля не прошли валидацию",
				TraceId:     traceId,
				Timestamp:   timestamp,
				Path:        path,
				FieldErrors: []errs.FieldError{fieldError},
			}
		}
	}

	logger.Error.Printf("[response.HandleError] traceId=%s Unexpected error: %s", traceId, err)
	return http.StatusInternalServerError, errs.ApiErrorResponse{
		Code:      errs.ErrorCodeInternalServerError,
		Message:   "Внутренняя ошибка сервера",
		TraceId:   traceId,
		Timestamp: timestamp,
		Path:      path,
	}
}

func handleValidationError(err error, traceId, timestamp, path string) (int, interface{}) {
	var fieldErrors []errs.FieldError

	if validationErrs, ok := err.(*errs.ValidationErrors); ok {
		fieldErrors = validationErrs.FieldErrors
	} else if multiErr, ok := err.(*errs.MultiValidationError); ok {
		fieldErrors = multiErr.FieldErrors
	} else if singleErr, ok := err.(*errs.ValidationError); ok {
		fieldErrors = append(fieldErrors, errs.FieldError{
			Field: singleErr.Field,
			Issue: singleErr.Issue,
		})
	} else {
		for validationErr, fieldError := range validationErrorMappings {
			if errors.Is(err, validationErr) {
				fieldErrorCopy := fieldError
				fieldErrors = append(fieldErrors, fieldErrorCopy)
			}
		}

		if field, issue, ok := errs.GetValidationErrorDetails(err); ok {
			fieldErrors = append(fieldErrors, errs.FieldError{
				Field: field,
				Issue: issue,
			})
		}
	}

	if len(fieldErrors) == 0 {
		fieldErrors = append(fieldErrors, errs.FieldError{
			Field: "general",
			Issue: "Ошибка валидации",
		})
	}

	return http.StatusUnprocessableEntity, errs.ValidationErrorResponse{
		Code:        errs.ErrorCodeValidationFailed,
		Message:     "Некоторые поля не прошли валидацию",
		TraceId:     traceId,
		Timestamp:   timestamp,
		Path:        path,
		FieldErrors: fieldErrors,
	}
}

func HandleError(c *gin.Context, err error) {
	traceId := observability.GetTraceID(c.Request.Context())
	timestamp := date_time.GetCurrentTimestamp()
	path := getPathFromContext(c)

	statusCode, response := mapErrorToResponse(err, traceId, timestamp, path)
	c.AbortWithStatusJSON(statusCode, response)
}
