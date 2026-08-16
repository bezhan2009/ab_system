package errs

import "errors"

var (
	ErrActiveExperimentExists           = errors.New("активный эксперимент на флаг уже существует")
	ErrCannotArchiveNotCompleted        = errors.New("можно архивировать только завершенные эксперименты")
	ErrExperimentStatusIsNotOnReview    = errors.New("эксперимент не находится на ревю")
	ErrExperimentAlreadyHasBeenOnReview = errors.New("эксперимент уже ранее был на ревью или находится")
	ErrInvalidStatusTransition          = errors.New("статусы не могут так двигаться, неправильный жизненный цикл эксперимента")
	ErrExperimentNotApproved            = errors.New("эксперимент ещё не одобрен")
	ErrConclusionRequired               = errors.New("причина ОБЯЗАТЕЛЬНА при завершении эксперимента")
	ErrCannotEditRunningExperiment      = errors.New("нельзя редактировать запущенный или приостановленный эксперимент")
	ErrInvalidWinnerVariant             = errors.New("не найден выигрышный вариант или он не валиден")
)
