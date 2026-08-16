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

type EventTypeHandler struct {
	service *service.EventTypeService
}

func NewEventTypeHandler(service *service.EventTypeService) *EventTypeHandler {
	return &EventTypeHandler{service: service}
}

func (h *EventTypeHandler) GetAllEventTypes(c *gin.Context) {
	types, err := h.service.GetAllEventTypes(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)

		return
	}

	dtoObj := dto.EventType{}
	c.JSON(http.StatusOK, dtoObj.ToDTOs(types))
}

func (h *EventTypeHandler) GetEventTypeByID(c *gin.Context) {
	id := c.Param("id")

	et, err := h.service.GetEventTypeByID(c.Request.Context(), id)
	if err != nil {
		response.HandleError(c, err)

		return
	}

	dtoObj := dto.EventType{}
	c.JSON(http.StatusOK, dtoObj.ToDTO(et))
}

func (h *EventTypeHandler) CreateEventType(c *gin.Context) {
	var req dto.EventType
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)

		return
	}

	if err := validation.ValidateEventTypeCreate(req); err != nil {
		response.HandleError(c, errs.NewMultiValidationError(err))

		return
	}

	model, err := req.ToModel()
	if err != nil {
		response.HandleError(c, errs.ErrIdIsInvalid)

		return
	}

	if err = h.service.CreateEventType(c.Request.Context(), model); err != nil {
		response.HandleError(c, err)

		return
	}

	c.JSON(http.StatusCreated, req.ToDTO(model))
}

func (h *EventTypeHandler) UpdateEventType(c *gin.Context) {
	id := c.Param("id")

	var req dto.EventType

	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	req.ID = id

	model, err := req.ToModel()
	if err != nil {
		response.HandleError(c, errs.ErrIdIsInvalid)

		return
	}

	if err := validation.ValidateEventTypeUpdate(req); err != nil {
		response.HandleError(c, errs.NewMultiValidationError(err))

		return
	}

	updated, err := h.service.UpdateEventType(c.Request.Context(), model)
	if err != nil {
		response.HandleError(c, err)

		return
	}

	c.JSON(http.StatusOK, req.ToDTO(updated))
}

func (h *EventTypeHandler) DeleteEventType(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteEventType(c.Request.Context(), id); err != nil {
		response.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
