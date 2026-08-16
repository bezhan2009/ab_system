package handlers

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/service"
	"ab_system/internal/http/dto"
	"ab_system/internal/http/response"
	"ab_system/internal/validation"
	"ab_system/pkg/errs"
	"net/http"

	contextKeys "ab_system/internal/http"

	"github.com/gin-gonic/gin"
)

type ExperimentHandler struct {
	service service.ExperimentService
}

func NewExperimentHandler(service service.ExperimentService) *ExperimentHandler {
	return &ExperimentHandler{
		service: service,
	}
}

func (h *ExperimentHandler) GetAllExperiments(c *gin.Context) {
	title := c.Query("title")

	var exps *[]models.Experiment
	var err error

	if title == "" {
		exps, err = h.service.GetAllExperiments(c.Request.Context())
		if err != nil {
			response.HandleError(c, err)
			return
		}
	} else {
		exps, err = h.service.GetExperimentByTitleLike(c.Request.Context(), title)
		if err != nil {
			response.HandleError(c, err)
			return
		}
	}

	expDto := dto.Experiment{}

	c.JSON(http.StatusOK, expDto.ToDTOs(*exps))
}

func (h *ExperimentHandler) GetAllExperimentsByStatus(c *gin.Context) {
	status := c.Param("status")

	var exps *[]models.Experiment
	var err error

	if status == "" {
		exps, err = h.service.GetAllExperiments(c.Request.Context())
		if err != nil {
			response.HandleError(c, err)

			return
		}
	} else {
		exps, err = h.service.GetExperimentsByStatus(c.Request.Context(), status)
		if err != nil {
			response.HandleError(c, err)

			return
		}
	}

	expDto := dto.Experiment{}

	c.JSON(http.StatusOK, expDto.ToDTOs(*exps))
}

func (h *ExperimentHandler) GetExperimentByID(c *gin.Context) {
	expId := c.Param("id")

	exp, err := h.service.GetExperimentByID(c.Request.Context(), expId)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	expDto := dto.Experiment{}
	c.JSON(http.StatusOK, expDto.ToDTO(exp))
}

func (h *ExperimentHandler) GetExperimentVersionsByID(c *gin.Context) {
	expId := c.Param("id")

	expVersions, err := h.service.GetExperimentVersionsByID(c.Request.Context(), expId)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	expDto := dto.ExperimentVersion{}
	c.JSON(http.StatusOK, expDto.ToDTOs(*expVersions))
}

func (h *ExperimentHandler) GetRunningExperimentByFlag(c *gin.Context) {
	flag := c.Param("flag")
	status := c.Param("status")

	exps, err := h.service.GetExperimentsByFlagAndStatus(c.Request.Context(), flag, status)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	expDto := dto.Experiment{}
	c.JSON(http.StatusOK, expDto.ToDTOs(*exps))
}

func (h *ExperimentHandler) CreateExperiment(c *gin.Context) {
	userID := c.GetString(contextKeys.UserIDCtx)

	var expDTO dto.Experiment
	if err := c.ShouldBindJSON(&expDTO); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)

		return
	}

	expDTO.ID = ""
	expDTO.OwnerID = userID
	expDTO.Status = "draft"

	if err := validation.ValidateExperimentCreate(&expDTO); err != nil {
		response.HandleError(c, errs.NewMultiValidationError(err))

		return
	}

	expModel, err := expDTO.ToModel()
	if err != nil {
		response.HandleError(c, errs.ErrIdIsInvalid)

		return
	}

	err = h.service.CreateExperiment(c.Request.Context(), expModel)
	if err != nil {
		response.HandleError(c, err)

		return
	}

	c.JSON(http.StatusCreated, expDTO.ToDTO(expModel))
}

func (h *ExperimentHandler) UpdateExperiment(c *gin.Context) {
	experimentId := c.Param("id")
	userID := c.GetString(contextKeys.UserIDCtx)

	var expDTO dto.Experiment
	if err := c.ShouldBindJSON(&expDTO); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	expDTO.OwnerID = userID
	expDTO.ID = experimentId

	expModel, err := expDTO.ToModel()
	if err != nil {
		response.HandleError(c, errs.ErrIdIsInvalid)
		return
	}

	if err := validation.ValidateExperimentUpdate(&expDTO); err != nil {
		response.HandleError(c, errs.NewMultiValidationError(err))
		return
	}

	updated, err := h.service.UpdateExperiment(c.Request.Context(), expModel)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, expDTO.ToDTO(updated))
}

func (h *ExperimentHandler) SendExperimentToReview(c *gin.Context) {
	experimentId := c.Param("id")

	err := h.service.SendToReview(c.Request.Context(), experimentId)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully sent to review",
	})
}

func (h *ExperimentHandler) RunExperiment(c *gin.Context) {
	experimentId := c.Param("id")

	err := h.service.RunExperiment(c.Request.Context(), experimentId)
	if err != nil {
		response.HandleError(c, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully run experiment",
	})
}

func (h *ExperimentHandler) CompleteExperiment(c *gin.Context) {
	experimentId := c.Param("id")

	var req dto.CompleteExperimentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)

		return
	}

	if err := validation.ValidateCompleteExperiment(&req); err != nil {
		response.HandleError(c, errs.NewMultiValidationError(err))
		return
	}

	err := h.service.CompleteExperiment(c.Request.Context(), req, experimentId)
	if err != nil {
		response.HandleError(c, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully completed experiment",
	})
}

func (h *ExperimentHandler) ArchiveExperiment(c *gin.Context) {
	experimentId := c.Param("id")

	err := h.service.ArchiveExperiment(c.Request.Context(), experimentId)
	if err != nil {
		response.HandleError(c, err)

		return
	}

	c.Status(http.StatusNoContent)
}
