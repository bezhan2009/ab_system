package handlers

import (
	"ab_system/internal/domain/service"
	contextkeysHttp "ab_system/internal/http"
	"ab_system/internal/http/dto"
	"ab_system/internal/http/response"
	"ab_system/internal/validation"
	"ab_system/pkg/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApprovalHandlers struct {
	approvalService service.ApprovalService
}

func NewApprovalHandlers(approvalService service.ApprovalService) *ApprovalHandlers {
	return &ApprovalHandlers{approvalService: approvalService}
}

func (h *ApprovalHandlers) GetApprovalsByExperimentID(c *gin.Context) {
	experimentID := c.Param("id")

	approvals, err := h.approvalService.GetApprovalsByExperimentID(c.Request.Context(), experimentID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	approvalsDto := dto.Approval{}

	c.JSON(http.StatusOK, dto.OrdinaryResponse[*dto.Approval]{
		Items: approvalsDto.ToDTOs(*approvals),
		Total: int64(len(*approvals)),
	})
}

func (h *ApprovalHandlers) GetApprovalByID(c *gin.Context) {
	approvalId := c.Param("id")

	approval, err := h.approvalService.GetApprovalByID(c.Request.Context(), approvalId)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	approvalDTO := dto.Approval{}
	c.JSON(http.StatusOK, approvalDTO.ToDTO(approval))
}

func (h *ApprovalHandlers) CreateApproval(c *gin.Context) {
	userID := c.GetString(contextkeysHttp.UserIDCtx)

	var approval dto.Approval
	if err := c.ShouldBindJSON(&approval); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)

		return
	}

	approval.ApproverID = userID

	if err := validation.ValidateApprovalCreate(&approval); err != nil {
		response.HandleError(c, errs.NewMultiValidationError(err))

		return
	}

	approvalModel, err := approval.ToModel()
	if err != nil {
		response.HandleError(c, errs.ErrIdIsInvalid)

		return
	}

	err = h.approvalService.CreateApproval(c.Request.Context(), approvalModel)
	if err != nil {
		response.HandleError(c, err)

		return
	}

	c.JSON(http.StatusOK, approval.ToDTO(approvalModel))
}

func (h *ApprovalHandlers) DeleteApproval(c *gin.Context) {
	approvalId := c.Param("id")

	err := h.approvalService.DeleteApproval(c.Request.Context(), approvalId)
	if err != nil {
		response.HandleError(c, err)

		return
	}

	c.Status(http.StatusNoContent)
}
