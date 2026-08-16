package validation

import (
	"ab_system/internal/http/dto"
	"ab_system/pkg/errs"
)

func ValidateApprovalCreate(approval *dto.Approval) (err []errs.FieldError) {
	if approval.ExperimentID == "" {
		err = append(err, errs.NewFieldError("experiment_id", "Обязательное поле", approval.ExperimentID))
	}

	if approval.ApproverID == "" {
		err = append(err, errs.NewFieldError("approver_id", "Обязательное поле", approval.ApproverID))
	}

	return err
}
