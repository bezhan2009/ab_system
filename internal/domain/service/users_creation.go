package service

import (
	"ab_system/internal/configs"
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	jwtauth "ab_system/internal/lib/jwt"
	"ab_system/pkg/errs"
	"ab_system/pkg/utils"
	"context"
	"errors"
)

type UserCreator struct {
	userWriter repository.UserWriter
	userReader repository.UserReader
	teamReader repository.TeamReader
	cfg        configs.Configs
}

func NewUserCreator(
	userWriter repository.UserWriter,
	userReader repository.UserReader,
) *UserCreator {
	return &UserCreator{
		userWriter: userWriter,
		userReader: userReader,
	}
}

func (uc *UserCreator) CreateUser(
	ctx context.Context,
	user *models.User,
	opts CreateUserOptions,
) (*models.User, string, error) {

	_, err := uc.userReader.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return nil, "", errs.ErrEmailUniquenessFailed
	}
	if !errors.Is(err, errs.ErrRecordNotFound) {
		return nil, "", err
	}

	if user.TeamID != nil {
		_, err = uc.teamReader.GetTeamByID(ctx, user.TeamID.String())
		if err != nil {
			return nil, "", err
		}
	}

	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		return nil, "", err
	}

	user.Role = opts.Role
	user.IsActive = true

	userDB, err := uc.userWriter.CreateUser(ctx, user)
	if err != nil {
		return nil, "", err
	}

	if !opts.GenerateToken {
		return userDB, "", nil
	}

	token, err := jwtauth.NewToken(
		userDB.ID.String(),
		string(userDB.Role),
		uc.cfg.AppParams.ServerName,
	)
	if err != nil {
		return nil, "", err
	}

	return userDB, token, nil
}
