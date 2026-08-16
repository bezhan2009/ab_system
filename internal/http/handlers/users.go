package handlers

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/service"
	contextkeysHttp "ab_system/internal/http"
	"ab_system/internal/http/dto"
	"ab_system/internal/http/response"
	"ab_system/internal/validation"
	"ab_system/pkg/errs"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
	userCreator service.UserCreator
}

func NewUserHandler(userService service.UserService, userCreator service.UserCreator) *UserHandler {
	return &UserHandler{userService: userService, userCreator: userCreator}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	var page int
	var size int
	var err error

	pageStr := c.DefaultQuery("page", "0")
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 0 {
			response.HandleError(c, errs.NewValidationError("page", "должен быть числом ≥ 0"))
			return
		}
	} else {
		page = 0
	}

	sizeStr := c.DefaultQuery("size", "20")
	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err != nil || size < 1 || size > 100 {
			response.HandleError(c, errs.NewValidationError("size", "должен быть числом от 1 до 100"))
			return
		}
	} else {
		size = 20
	}

	users, total, err := h.userService.GetAllUsers(
		c.Request.Context(),
		page,
		size,
	)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	userDTO := dto.User{}
	items := userDTO.ToDTOs(users)
	if items == nil {
		items = []*dto.User{}
	}

	c.JSON(http.StatusOK, dto.PageResponse[*dto.User]{
		Items: items,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := c.GetString(contextkeysHttp.UserIDCtx)
	currentUserRole := c.GetString(contextkeysHttp.UserRoleCtx)

	if currentUserRole == "USER" && userID != currentUserID {
		response.HandleError(c, errs.ErrPermissionDenied)
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	userDTO := dto.User{}
	c.JSON(http.StatusOK, userDTO.ToDTO(user))
}

func (h *UserHandler) GetUserMe(c *gin.Context) {
	userID := c.GetString(contextkeysHttp.UserIDCtx)

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	userDTO := dto.User{}
	c.JSON(http.StatusOK, userDTO.ToDTO(user))
}

func (h *UserHandler) GetUsersByTeamID(c *gin.Context) {
	teamID := c.Param("id")

	users, err := h.userService.GetTeamUsersByTeamID(c.Request.Context(), teamID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	userDTO := dto.User{}
	items := userDTO.ToDTOs(*users)
	if items == nil {
		items = []*dto.User{}
	}

	c.JSON(http.StatusOK, dto.OrdinaryResponse[*dto.User]{
		Items: items,
		Total: int64(len(items)),
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var createReq dto.CreateUserReq

	if err := c.ShouldBindJSON(&createReq); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	var validationErrors []errs.FieldError

	if createReq.Role == "" {
		validationErrors = append(validationErrors,
			errs.NewFieldError("role", "обязательное поле"))
	}

	user := dto.User{
		Email:    createReq.Email,
		Password: createReq.Password,
		FullName: createReq.FullName,
		Region:   createReq.Region,
		Gender:   createReq.Gender,
		Age:      createReq.Age,
		Role:     createReq.Role,
		IsActive: true,
	}

	if valueErrors := validation.ValidateCreateUser(&user); len(valueErrors) > 0 {
		validationErrors = append(validationErrors, valueErrors...)
	}

	if len(validationErrors) > 0 {
		response.HandleError(c, errs.NewValidationErrors(validationErrors))
		return
	}

	userDB, _, err := h.userCreator.CreateUser(
		c.Request.Context(),
		user.ToModel(),
		service.CreateUserOptions{
			Role:          models.Role(user.Role),
			GenerateToken: false,
		},
	)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	userDTO := dto.User{}
	c.JSON(http.StatusCreated, userDTO.ToDTO(userDB))
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	currentUserID := c.GetString(contextkeysHttp.UserIDCtx)
	currentUserRole := c.GetString(contextkeysHttp.UserRoleCtx)

	if currentUserRole == "USER" && userID != currentUserID {
		response.HandleError(c, errs.ErrPermissionDenied)
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	requiredKeys := []string{"fullName", "age", "region", "gender", "maritalStatus"}
	var missingFields []errs.FieldError
	for _, key := range requiredKeys {
		if _, exists := raw[key]; !exists {
			missingFields = append(missingFields, errs.NewFieldError(key, "обязательное поле"))
		}
	}

	if len(missingFields) > 0 {
		response.HandleError(c, errs.NewValidationErrors(missingFields))
		return
	}

	var updateReq struct {
		FullName      *string            `json:"fullName"`
		Age           *int               `json:"age"`
		Region        *string            `json:"region"`
		Gender        *dto.Gender        `json:"gender"`
		Role          *dto.Role          `json:"role"`
		IsActive      *bool              `json:"isActive"`
	}

	if err := json.Unmarshal(body, &updateReq); err != nil {
		response.HandleError(c, errs.ErrValidationFailed)
		return
	}

	if updateReq.FullName == nil {
		response.HandleError(c, errs.NewValidationErrors([]errs.FieldError{
			errs.NewFieldError("fullName", "не может быть null"),
		}))
		return
	}

	if currentUserRole != "admin" {
		if updateReq.Role != nil {
			response.HandleError(c, errs.ErrPermissionDenied)
			return
		}
		if updateReq.IsActive != nil {
			response.HandleError(c, errs.ErrPermissionDenied)
			return
		}
	}

	user := dto.User{
		ID:       userID,
		FullName: *updateReq.FullName,
	}

	if updateReq.Age != nil {
		user.Age = updateReq.Age
	} else {
		user.Age = nil
	}

	if updateReq.Region != nil {
		user.Region = updateReq.Region
	} else {
		user.Region = nil
	}

	if updateReq.Gender != nil {
		user.Gender = updateReq.Gender
	} else {
		user.Gender = nil
	}

	if currentUserRole == "admin" {
		if updateReq.Role != nil {
			user.Role = *updateReq.Role
		}
		if updateReq.IsActive != nil {
			user.IsActive = *updateReq.IsActive
		}
	} else {
		existingUser, err := h.userService.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			response.HandleError(c, err)
			return
		}
		user.Role = dto.Role(existingUser.Role)
		user.IsActive = existingUser.IsActive
	}

	if valueErrors := validation.ValidateUpdateUserFull(&user); len(valueErrors) > 0 {
		response.HandleError(c, errs.NewValidationErrors(valueErrors))
		return
	}

	userDB, err := h.userService.UpdateUser(
		c.Request.Context(),
		user.ToModel(),
	)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	userDTO := dto.User{}
	c.JSON(http.StatusOK, userDTO.ToDTO(userDB))
}

func (h *UserHandler) UpdateUserMe(c *gin.Context) {
	userID := c.GetString(contextkeysHttp.UserIDCtx)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: userID})
	h.UpdateUser(c)
}

func (h *UserHandler) DeactivateUser(c *gin.Context) {
	userID := c.Param("id")

	err := h.userService.DeactivateUser(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
