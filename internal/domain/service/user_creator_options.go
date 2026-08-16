package service

import "ab_system/internal/domain/models"

type CreateUserOptions struct {
	Role          models.Role
	GenerateToken bool
}
