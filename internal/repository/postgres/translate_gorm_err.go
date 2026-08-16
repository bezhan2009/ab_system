package postgres

import (
	"ab_system/pkg/errs"
	"errors"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

func TranslateGormError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.ErrRecordNotFound
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return errs.ErrDuplicateEntry
	}
	if errors.Is(err, gorm.ErrInvalidField) {
		return errs.ErrInvalidField
	}
	if errors.Is(err, gorm.ErrInvalidData) {
		return errs.ErrInvalidData
	}

	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return errs.ErrDuplicateEntry
		case "23503":
			return errs.ErrInvalidField
		case "23502":
			return errs.ErrInvalidData
		case "22001":
			return errs.ErrInvalidData
		}
	}

	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "duplicate key") ||
		strings.Contains(errStr, "unique constraint") ||
		strings.Contains(errStr, "повторяющееся") {
		return errs.ErrDuplicateEntry
	}
	if strings.Contains(errStr, "foreign key") {
		return errs.ErrInvalidField
	}

	return err
}

func isForeignKeyViolation(err error) bool {
	return strings.Contains(err.Error(), "violates foreign key constraint")
}
