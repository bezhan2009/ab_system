package handlers

import (
	"ab_system/internal/domain/service"
	"ab_system/internal/http/response"
	"ab_system/pkg/errs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	service service.ReportService
}

func NewReportHandler(service service.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) GetExperimentReport(c *gin.Context) {
	expID := c.Param("id")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	var from, to time.Time
	var err error

	if toStr == "" {
		now := time.Now()
		to = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	} else {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			response.HandleError(c, errs.ErrInvalidDateFormat)
			return
		}
	}

	if fromStr == "" {
		from = to.AddDate(0, -1, 0)
	} else {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			response.HandleError(c, errs.ErrInvalidDateFormat)
			return
		}
	}

	report, err := h.service.GetExperimentReport(c.Request.Context(), expID, from, to)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, report)
}
