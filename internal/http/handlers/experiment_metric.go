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

type ExperimentMetricHandler struct {
	service *service.ExperimentMetricService
}

func NewExperimentMetricHandler(service *service.ExperimentMetricService) *ExperimentMetricHandler {
	return &ExperimentMetricHandler{service: service}
}

func (h *ExperimentMetricHandler) GetExperimentMetrics(c *gin.Context) {
	expID := c.Param("id")

	metrics, err := h.service.GetExperimentMetrics(c.Request.Context(), expID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	dtoObj := dto.ExperimentMetric{}
	c.JSON(http.StatusOK, dtoObj.ToDTOs(metrics))
}

func (h *ExperimentMetricHandler) GetGuardrailsForExperiment(c *gin.Context) {
	expID := c.Param("id")

	guardrails, err := h.service.GetGuardrailsForExperiment(c.Request.Context(), expID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	dtoObj := dto.ExperimentMetric{}
	c.JSON(http.StatusOK, dtoObj.ToDTOs(guardrails))
}

func (h *ExperimentMetricHandler) AddMetricToExperiment(c *gin.Context) {
	expID := c.Param("id")
	var req dto.ExperimentMetric

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

	if model.IsGuardrail {
		if err := validation.ValidateGuardrailConfig(model); err != nil {
			response.HandleError(c, errs.NewMultiValidationError(err))

			return
		}
	}

	if err = h.service.AddMetricToExperiment(c.Request.Context(), model); err != nil {
		response.HandleError(c, err)

		return
	}

	c.JSON(http.StatusCreated, req.ToDTO(model))
}

func (h *ExperimentMetricHandler) UpdateExperimentMetric(c *gin.Context) {
	expID := c.Param("id")
	metricID := c.Param("metric_id")
	var req dto.ExperimentMetric

	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	req.ExperimentID = expID
	req.MetricID = metricID

	model, err := req.ToModel()
	if err != nil {
		response.HandleError(c, errs.ErrIdIsInvalid)

		return
	}

	if model.IsGuardrail {
		if err := validation.ValidateGuardrailConfig(model); err != nil {
			response.HandleError(c, errs.NewMultiValidationError(err))

			return
		}
	}

	updated, err := h.service.UpdateExperimentMetric(c.Request.Context(), model)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, req.ToDTO(updated))
}

func (h *ExperimentMetricHandler) RemoveMetricFromExperiment(c *gin.Context) {
	expID := c.Param("id")
	metricID := c.Param("metric_id")

	err := h.service.RemoveMetricFromExperiment(c.Request.Context(), expID, metricID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
