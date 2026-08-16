package handlers

import (
	"ab_system/internal/domain/service"
	"ab_system/internal/http/dto"
	"ab_system/internal/http/response"
	"ab_system/pkg/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApproverGroupHandler struct {
	service *service.ApproverGroupService
}

func NewApproverGroupHandler(service *service.ApproverGroupService) *ApproverGroupHandler {
	return &ApproverGroupHandler{service: service}
}

func (h *ApproverGroupHandler) GetApproverGroupByExperimenterID(c *gin.Context) {
	experimenterID := c.Param("id")

	group, members, err := h.service.GetGroupByExperimenterID(c.Request.Context(), experimenterID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	dtoObj := dto.ApproverGroup{}
	c.JSON(http.StatusOK, dtoObj.ToDTO(group, members))
}

func (h *ApproverGroupHandler) GetApproverGroupByExperimentID(c *gin.Context) {
	experimentID := c.Param("id")

	group, members, err := h.service.GetApproverGroupByExperimentID(c.Request.Context(), experimentID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	dtoObj := dto.ApproverGroup{}
	c.JSON(http.StatusOK, dtoObj.ToDTO(group, members))
}

func (h *ApproverGroupHandler) CreateApproverGroup(c *gin.Context) {
	var req dto.ApproverGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	model, approverUUIDs, err := req.ToModel()
	if err != nil {
		response.HandleError(c, errs.ErrIdIsInvalid)
		return
	}

	approverIDs := make([]string, len(approverUUIDs))
	for i, uid := range approverUUIDs {
		approverIDs[i] = uid.String()
	}

	if err = h.service.CreateApproverGroup(c.Request.Context(), model, approverIDs); err != nil {
		response.HandleError(c, err)
		return
	}

	createdGroup, members, err := h.service.GetGroupByExperimenterID(c.Request.Context(), model.ExperimenterID.String())
	if err != nil {
		c.JSON(http.StatusMultiStatus, gin.H{"id": model.ID.String()})
		return
	}

	dtoObj := dto.ApproverGroup{}
	c.JSON(http.StatusCreated, dtoObj.ToDTO(createdGroup, members))
}

func (h *ApproverGroupHandler) UpdateApproverGroup(c *gin.Context) {
	id := c.Param("id")

	var req dto.ApproverGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	req.ID = id

	model, approverUUIDs, err := req.ToModel()
	if err != nil {
		response.HandleError(c, errs.ErrIdIsInvalid)
		return
	}

	approverIDs := make([]string, len(approverUUIDs))
	for i, uid := range approverUUIDs {
		approverIDs[i] = uid.String()
	}

	updated, err := h.service.UpdateApproverGroup(c.Request.Context(), model, approverIDs)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	_, members, err := h.service.GetGroupByExperimenterID(c.Request.Context(), updated.ExperimenterID.String())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"id": updated.ID.String()})
		return
	}

	dtoObj := dto.ApproverGroup{}
	c.JSON(http.StatusOK, dtoObj.ToDTO(updated, members))
}

func (h *ApproverGroupHandler) DeleteApproverGroup(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteApproverGroup(c.Request.Context(), id); err != nil {
		response.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
