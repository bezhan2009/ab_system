package errs

type NotFoundWithDetailsError struct {
	Message string
	Details map[string]interface{}
}

func (e *NotFoundWithDetailsError) Error() string {
	return e.Message
}
