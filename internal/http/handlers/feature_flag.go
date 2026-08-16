package handlers

import (
	"ab_system/internal/domain/service"
	"ab_system/internal/http/dto"
	"ab_system/internal/http/response"
	"ab_system/internal/validation"
	"ab_system/pkg/errs"
	"net/http"

	contextKeys "ab_system/internal/http"

	"github.com/gin-gonic/gin"
)

type FeatureFlagHandler struct {
	service service.FeatureFlagService
}

func NewFeatureFlagHandler(service service.FeatureFlagService) *FeatureFlagHandler {
	return &FeatureFlagHandler{
		service: service,
	}
}

func (h *FeatureFlagHandler) GetAllFeatureFlags(c *gin.Context) {
	featureFlags, err := h.service.GetAllFeatureFlags(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}

	featureFlagsDto := dto.FeatureFlag{}

	c.JSON(http.StatusOK, featureFlagsDto.ToDTOs(featureFlags))
}

func (h *FeatureFlagHandler) GetFeatureFlagById(c *gin.Context) {
	featureFlagId := c.Param("id")

	featureFlag, err := h.service.GetFeatureFlagById(c.Request.Context(), featureFlagId)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	featureFlagDto := dto.FeatureFlag{}

	c.JSON(http.StatusOK, featureFlagDto.ToDTO(&featureFlag))
}

func (h *FeatureFlagHandler) GetFeatureFlagsByKey(c *gin.Context) {
	featureFlagKey := c.Param("key")

	featureFlag, err := h.service.GetFeatureFlagsByKey(c.Request.Context(), featureFlagKey)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	featureFlagDto := dto.FeatureFlag{}

	c.JSON(http.StatusOK, featureFlagDto.ToDTO(&featureFlag))
}

func (h *FeatureFlagHandler) GetFeatureFlagsByOwner(c *gin.Context) {
	owner := c.Param("id")

	featureFlags, err := h.service.GetFeatureFlagsByOwner(c.Request.Context(), owner)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	featureFlagDto := dto.FeatureFlag{}

	c.JSON(http.StatusOK, featureFlagDto.ToDTOs(featureFlags))
}

func (h *FeatureFlagHandler) CreateFeatureFlag(c *gin.Context) {
	userId := c.GetString(contextKeys.UserIDCtx)

	var featureFlag dto.FeatureFlag

	if err := c.ShouldBindJSON(&featureFlag); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	featureFlag.UserID = userId

	if err := validation.ValidateFeatureFlagCreate(featureFlag); err != nil {
		response.HandleError(c, errs.NewMultiValidationError(err))

		return
	}

	featureFlagModel, err := featureFlag.ToModel()
	if err != nil {
		response.HandleError(c, errs.ErrIdIsInvalid)

		return
	}

	err = h.service.CreateFeatureFlag(c.Request.Context(), featureFlagModel)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, featureFlag.ToDTO(featureFlagModel))
}

func (h *FeatureFlagHandler) UpdateFeatureFlag(c *gin.Context) {
	flagId := c.Param("id")
	userId := c.GetString(contextKeys.UserIDCtx)

	var featureFlag dto.FeatureFlag

	if err := c.ShouldBindJSON(&featureFlag); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	featureFlag.UserID = userId
	featureFlag.ID = flagId

	featureFlagModel, err := featureFlag.ToModel()
	if err != nil {
		response.HandleError(c, errs.ErrIdIsInvalid)
		return
	}

	updatedFeatureFlag, err := h.service.UpdateFeatureFlag(c.Request.Context(), featureFlagModel)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, featureFlag.ToDTO(updatedFeatureFlag))
}

func (h *FeatureFlagHandler) DeleteFeatureFlagById(c *gin.Context) {
	featureFlagId := c.Param("id")

	err := h.service.DeleteFeatureFlagById(c.Request.Context(), featureFlagId)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
