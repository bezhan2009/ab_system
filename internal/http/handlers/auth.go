package handlers

import (
	"ab_system/internal/configs"
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/service"
	"ab_system/internal/http/dto"
	"ab_system/internal/http/response"
	"ab_system/internal/validation"
	"ab_system/pkg/errs"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth        service.AuthService
	userCreator service.UserCreator
	cfg         configs.Configs
}

func NewAuthHandler(service service.AuthService, userCreator service.UserCreator, cfg configs.Configs) *AuthHandler {
	return &AuthHandler{auth: service, userCreator: userCreator, cfg: cfg}
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	var auth dto.Auth
	if err := c.ShouldBind(&auth); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	if err := validation.ValidateLogin(auth.Email, auth.Password); err != nil {
		response.HandleError(c, errs.NewMultiValidationError(err))
		return
	}

	user, accessToken, err := h.auth.Login(c.Request.Context(), auth.Email, auth.Password)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	userDto := dto.User{}

	jwtTtlSecondsInt, _ := strconv.ParseInt(os.Getenv("JWT_TTL_SECONDS"), 10, 64)

	c.JSON(http.StatusOK, dto.AuthResponse{
		AccessToken: accessToken,
		ExpiresIn:   jwtTtlSecondsInt,
		User:        *userDto.ToDTO(user),
	})
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var user dto.CreateUserReq
	if err := c.ShouldBind(&user); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	user.IsActive = true
	user.Role = dto.Role(models.RoleViewer)

	userDto := dto.User{
		Email:    user.Email,
		FullName: user.FullName,
		Password: user.Password,
		Region:   user.Region,
		Gender:   user.Gender,
		Age:      user.Age,
		Role:     user.Role,
		IsActive: user.IsActive,
	}

	if err := validation.ValidateCreateUser(&userDto); err != nil {
		response.HandleError(c, errs.NewMultiValidationError(err))
		return
	}

	userDB, accessToken, err := h.userCreator.CreateUser(
		c.Request.Context(),
		userDto.ToModel(),
		service.CreateUserOptions{
			Role:          models.Role(userDto.Role),
			GenerateToken: true,
		},
	)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	jwtTtlSecondsInt, _ := strconv.ParseInt(os.Getenv("JWT_TTL_SECONDS"), 10, 64)

	c.JSON(http.StatusCreated, dto.AuthResponse{
		AccessToken: accessToken,
		ExpiresIn:   jwtTtlSecondsInt,
		User:        *userDto.ToDTO(userDB),
	})
}
