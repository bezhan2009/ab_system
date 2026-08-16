package handlers

import (
	"ab_system/internal/domain/service"
	"ab_system/internal/http/dto"
	"ab_system/internal/http/response"
	"ab_system/pkg/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationSettingsHandler struct {
	service *service.NotificationSettingsService
}

func NewNotificationSettingsHandler(service *service.NotificationSettingsService) *NotificationSettingsHandler {
	return &NotificationSettingsHandler{service: service}
}

func (h *NotificationSettingsHandler) GetNotificationSettingsByExperimentID(c *gin.Context) {
	expID := c.Param("id")

	settings, err := h.service.GetNotificationSettingsByExperimentID(c.Request.Context(), expID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	dtoObj := dto.NotificationSettings{}
	resp, err := dtoObj.ToDTO(settings)
	if err != nil {
		response.HandleError(c, errs.ErrInvalidData)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *NotificationSettingsHandler) CreateNotificationSettings(c *gin.Context) {
	expID := c.Param("id")

	var req dto.NotificationSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	req.ExperimentID = expID

	model, err := req.ToModel()
	if err != nil {
		response.HandleError(c, errs.ErrIdIsInvalid)

		return
	}

	settings, err := h.service.CreateNotificationSettings(c.Request.Context(), model)
	if err != nil {
		response.HandleError(c, err)

		return
	}

	resp, err := req.ToDTO(settings)
	if err != nil {
		response.HandleError(c, errs.ErrInvalidData)

		return

	}

	c.JSON(http.StatusOK, resp)
}

func (h *NotificationSettingsHandler) UpdateNotificationSettings(c *gin.Context) {
	expID := c.Param("id")

	var req dto.NotificationSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)

		return
	}

	req.ExperimentID = expID

	model, err := req.ToModel()
	if err != nil {
		response.HandleError(c, errs.ErrIdIsInvalid)

		return
	}

	settings, err := h.service.UpdateNotificationSettings(c.Request.Context(), model)
	if err != nil {
		response.HandleError(c, err)

		return
	}

	resp, err := req.ToDTO(settings)
	if err != nil {
		response.HandleError(c, errs.ErrInvalidData)

		return
	}

	c.JSON(http.StatusOK, resp)
}

//func (h *NotificationSettingsHandler) DeleteNotificationSettings(c *gin.Context) {
//	expID := c.Param("id")
//
//	settings, err := h.service.GetNotificationSettingsByExperimentID(c.Request.Context(), expID)
//	if err != nil {
//		response.HandleError(c, err)
//
//		return
//	}
//
//	if err = h.service.DeleteNotificationSettingsByExperimentID(c.Request.Context(), settings.ID.String()); err != nil {
//		response.HandleError(c, err)
//
//		return
//	}
//
//	c.Status(http.StatusNoContent)
//}
