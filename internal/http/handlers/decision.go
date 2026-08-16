package handlers

import (
	"ab_system/internal/domain/service"
	"ab_system/internal/http/dto"
	"ab_system/internal/http/response"
	"ab_system/pkg/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DecisionHandler struct {
	decisionService service.DecisionService
}

func NewDecisionHandler(decisionService service.DecisionService) *DecisionHandler {
	return &DecisionHandler{decisionService: decisionService}
}

func (h *DecisionHandler) Decide(c *gin.Context) {
	var decisionReq dto.DecideRequest
	if err := c.ShouldBind(&decisionReq); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	decide, err := h.decisionService.Decide(c.Request.Context(), &decisionReq)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, decide)
}
