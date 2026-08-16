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

type MetricHandler struct {
	service *service.MetricService
}

func NewMetricHandler(service *service.MetricService) *MetricHandler {
	return &MetricHandler{service: service}
}

func (h *MetricHandler) GetAllMetrics(c *gin.Context) {
	metrics, err := h.service.GetAllMetrics(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}

	dtoObj := dto.Metric{}
	c.JSON(http.StatusOK, dtoObj.ToDTOs(metrics))
}

func (h *MetricHandler) GetMetricByID(c *gin.Context) {
	metricId := c.Param("id")

	metric, err := h.service.GetMetricByID(c.Request.Context(), metricId)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	dtoObj := dto.Metric{}
	c.JSON(http.StatusOK, dtoObj.ToDTO(metric))
}

func (h *MetricHandler) GetMetricByTitle(c *gin.Context) {
	title := c.Param("title")
	metric, err := h.service.GetMetricByTitle(c.Request.Context(), title)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	dtoObj := dto.Metric{}
	c.JSON(http.StatusOK, dtoObj.ToDTO(metric))
}

func (h *MetricHandler) CreateMetric(c *gin.Context) {
	var req dto.Metric
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	if errsList := validation.ValidateMetricCreate(req); len(errsList) > 0 {
		response.HandleError(c, errs.NewMultiValidationError(errsList))
		return
	}

	model, err := req.ToModel()
	if err != nil {
		response.HandleError(c, errs.ErrIdIsInvalid)
		return
	}

	if err = h.service.CreateMetric(c.Request.Context(), model); err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, req.ToDTO(model))
}

func (h *MetricHandler) UpdateMetric(c *gin.Context) {
	id := c.Param("id")
	var req dto.Metric
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	req.ID = id

	if errsList := validation.ValidateMetricUpdate(req); len(errsList) > 0 {
		response.HandleError(c, errs.NewMultiValidationError(errsList))
		return
	}

	model, err := req.ToModel()
	if err != nil {
		response.HandleError(c, errs.ErrIdIsInvalid)
		return
	}

	updated, err := h.service.UpdateMetric(c.Request.Context(), model)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, req.ToDTO(updated))
}

func (h *MetricHandler) DeleteMetric(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteMetric(c.Request.Context(), id); err != nil {
		response.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
