package service

import (
	"ab_system/internal/configs"
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	jwtauth "ab_system/internal/lib/jwt"
	"ab_system/pkg/errs"
	"ab_system/pkg/utils"
	"context"
)

type AuthService struct {
	userReader repository.UserReader
	cfg        configs.Configs
}

func NewAuthService(userReader repository.UserReader, cfg configs.Configs) *AuthService {
	return &AuthService{userReader: userReader, cfg: cfg}
}

func (auth *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (user *models.User, accessToken string, err error) {

	user, err = auth.userReader.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", errs.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, "", errs.ErrUserInactive
	}

	if !utils.CheckPassword(password, user.Password) {
		return nil, "", errs.ErrInvalidCredentials
	}

	accessToken, err = jwtauth.NewToken(
		user.ID.String(),
		string(user.Role),
		auth.cfg.AppParams.ServerName,
	)
	if err != nil {
		return nil, "", err
	}

	return user, accessToken, nil
}
