package repository

import (
	"ab_system/internal/domain/models"
	"context"
)

type UserReader interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
}

type UserTeamReader interface {
	GetUsersByTeamID(ctx context.Context, teamID string) (users *[]models.User, err error)
}

type UserWriter interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
}

type UserLister interface {
	GetAllUsers(
		ctx context.Context,
		page int,
		size int,
	) ([]models.User, int64, error)
}

type UserStatusUpdater interface {
	UpdateUserStatus(ctx context.Context, userID string, isActive bool) error
}

type UserStatusChecker interface {
	IsActive(ctx context.Context, userID string) (bool, error)
}
