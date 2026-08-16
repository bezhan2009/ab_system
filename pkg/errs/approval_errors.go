package errs

import "errors"

var (
	ErrForbiddenToBeApprover = errors.New("вы не имеете доступ на принятия экспериментов данного аналитика")
	ErrApproverAlreadyVoted  = errors.New("вы ранее уже отдавали голос за этот эксперимент")
)
