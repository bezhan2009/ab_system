package handlers

import (
	"ab_system/internal/domain/service"
	"ab_system/internal/http/dto"
	"ab_system/internal/http/response"
	"ab_system/pkg/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	service *service.TeamService
}

func NewTeamHandler(service *service.TeamService) *TeamHandler {
	return &TeamHandler{service: service}
}

func (h *TeamHandler) GetAllTeams(c *gin.Context) {
	name := c.Query("name")

	teams, err := h.service.GetAllTeams(c.Request.Context(), name)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	dtoObj := dto.Team{}
	c.JSON(http.StatusOK, dtoObj.ToDTOs(*teams))
}

func (h *TeamHandler) GetTeamByID(c *gin.Context) {
	id := c.Param("id")

	team, err := h.service.GetTeamByID(c.Request.Context(), id)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	dtoObj := dto.Team{}
	c.JSON(http.StatusOK, dtoObj.ToDTO(team))
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req dto.Team
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	model, err := req.ToModel()
	if err != nil {
		response.HandleError(c, errs.ErrIdIsInvalid)
		return
	}

	if err = h.service.CreateTeam(c.Request.Context(), model); err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, req.ToDTO(model))
}

func (h *TeamHandler) AddMemberToTeam(c *gin.Context) {
	var req dto.TeamMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	err := h.service.AddTeamMember(c.Request.Context(), req.TeamID, req.UserID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Пользователь успешно добавлен в команду",
	})
}

func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	id := c.Param("id")
	var req dto.Team
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

	updated, err := h.service.UpdateTeam(c.Request.Context(), model)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, req.ToDTO(updated))
}

func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteTeam(c.Request.Context(), id); err != nil {
		response.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
