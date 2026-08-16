package dto

import (
	"ab_system/internal/domain/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
)

type MaritalStatus string

const (
	MaritalSingle   MaritalStatus = "SINGLE"
	MaritalMarried  MaritalStatus = "MARRIED"
	MaritalDivorced MaritalStatus = "DIVORCED"
	MaritalWidowed  MaritalStatus = "WIDOWED"
)

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
)

type User struct {
	ID       string  `json:"id"`
	TeamID   *string `json:"team_id"`
	Email    string  `json:"email"`
	FullName string  `json:"fullName"`
	Password string  `json:"-"`
	Region   *string `json:"region"`
	Gender   *Gender `json:"gender"`
	Age      *int    `json:"age"`
	Role     Role    `json:"role"`

	IsActive  bool           `json:"isActive"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type RegisterRequest struct {
	Email    string  `json:"email" binding:"required,email,max=254"`
	Password string  `json:"password" binding:"required,min=8,max=72,password"`
	FullName string  `json:"fullName" binding:"required,min=2,max=200"`
	Region   *string `json:"region" binding:"omitempty,max=32"`
	TeamID   *string `json:"team_id"`
	Gender   *Gender `json:"gender" binding:"omitempty,oneof=MALE FEMALE"`
	Age      *int    `json:"age" binding:"omitempty,min=18,max=120"`
}

type CreateUserReq struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	FullName string  `json:"fullName"`
	Region   *string `json:"region"`
	Gender   *Gender `json:"gender"`
	Age      *int    `json:"age"`
	Role     Role    `json:"role"`
	IsActive bool    `json:"isActive"`
}

func (u *User) ToModel() *models.User {

	idUUID, _ := uuid.Parse(u.ID)
	var teamIDPtr *uuid.UUID
	if u.TeamID != nil {
		tid, err := uuid.Parse(*u.TeamID)
		if err == nil {
			teamIDPtr = &tid
		}
	}

	var regionPtr *string
	if u.Region != nil {
		regionPtr = u.Region
	}

	var agePtr *int
	if u.Age != nil {
		agePtr = u.Age
	}

	return &models.User{
		ID:        idUUID,
		Email:     u.Email,
		TeamID:    teamIDPtr,
		FullName:  u.FullName,
		Password:  u.Password,
		Region:    regionPtr,
		Age:       agePtr,
		Role:      models.Role(u.Role),
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) ToDTO(userModel *models.User) *User {
	var regionPtr *string
	if userModel.Region != nil {
		regionPtr = userModel.Region
	}

	var agePtr *int
	if userModel.Age != nil {
		agePtr = userModel.Age
	}

	teamUUID := ""
	if userModel.TeamID != nil {
		teamUUID = userModel.TeamID.String()
	}

	return &User{
		ID:        userModel.ID.String(),
		Email:     userModel.Email,
		TeamID:    &teamUUID,
		FullName:  userModel.FullName,
		Password:  userModel.Password,
		Region:    regionPtr,
		Age:       agePtr,
		Role:      Role(userModel.Role),
		IsActive:  userModel.IsActive,
		CreatedAt: userModel.CreatedAt,
		UpdatedAt: userModel.UpdatedAt,
	}
}

func (u *User) ToDTOs(users []models.User) []*User {
	result := make([]*User, len(users))
	for i := range users {
		result[i] = u.ToDTO(&users[i])
	}

	return result
}
