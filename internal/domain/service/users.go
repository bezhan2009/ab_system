package service

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	"ab_system/pkg/errs"
	"context"
	"errors"
)

type UserService struct {
	userWriter        repository.UserWriter
	userStatusUpdater repository.UserStatusUpdater
	userReader        repository.UserReader
	userTeamReader    repository.UserTeamReader
	userLister        repository.UserLister
}

func NewUserService(userWriter repository.UserWriter, userStatusUpdater repository.UserStatusUpdater, userReader repository.UserReader, userTeamReader repository.UserTeamReader, userLister repository.UserLister) *UserService {
	return &UserService{userWriter: userWriter, userStatusUpdater: userStatusUpdater, userReader: userReader, userTeamReader: userTeamReader, userLister: userLister}
}

func (s *UserService) GetAllUsers(
	ctx context.Context,
	page int,
	size int,
) (users []models.User, total int64, err error) {

	return s.userLister.GetAllUsers(ctx, page, size)
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (user *models.User, err error) {
	user, err = s.userReader.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetTeamUsersByTeamID(ctx context.Context, teamID string) (users *[]models.User, err error) {
	users, err = s.userTeamReader.GetUsersByTeamID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) (userDB *models.User, err error) {
	u, err := s.GetUserByID(ctx, user.ID.String())
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}

		return nil, err
	}

	user.Email = u.Email
	user.Password = u.Password

	userDB, err = s.userWriter.UpdateUser(ctx, user)
	if err != nil {
		return userDB, err
	}

	return userDB, nil
}

func (s *UserService) UpdateUserMe(ctx context.Context, user *models.User) (userDB *models.User, err error) {
	if user.Role == models.RoleAdmin || user.IsActive != true {
		return nil, errs.ErrPermissionDenied
	}

	u, err := s.GetUserByID(ctx, user.ID.String())
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}

		return nil, err
	}

	user.Email = u.Email
	user.Password = u.Password

	userDB, err = s.userWriter.UpdateUser(ctx, user)
	if err != nil {
		return userDB, err
	}

	return userDB, nil
}

func (s *UserService) DeactivateUser(ctx context.Context, userID string) (err error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return errs.ErrUserNotFound
		}
		return err
	}

	if !user.IsActive {
		return nil
	}

	user.IsActive = false
	err = s.userStatusUpdater.UpdateUserStatus(ctx, user.ID.String(), user.IsActive)
	if err != nil {
		return err
	}

	return nil
}
