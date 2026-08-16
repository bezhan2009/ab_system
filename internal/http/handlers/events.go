package handlers

import (
	"ab_system/internal/domain/service"
	"ab_system/internal/http/dto"
	"ab_system/internal/http/response"
	"ab_system/internal/validation"
	"ab_system/pkg/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) PostEvent(c *gin.Context) {
	var req dto.EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)

		return
	}

	if err := validation.ValidateEventCreate(req); err != nil {
		response.HandleError(c, errs.NewMultiValidationError(err))

		return
	}

	err := h.service.ProcessEvent(c.Request.Context(), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "event created successfully",
	})
}

func (h *EventHandler) PostEvents(c *gin.Context) {
	var req dto.EventBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)

		return
	}

	if err := validation.ValidateEventBatch(req); err != nil {
		response.HandleError(c, errs.NewMultiValidationError(err))

		return
	}

	resp, err := h.service.ProcessEvents(c.Request.Context(), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
