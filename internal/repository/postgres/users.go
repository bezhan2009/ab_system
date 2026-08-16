package postgres

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/http/middlewares/observability"
	"ab_system/pkg/errs"
	"ab_system/pkg/logger"
	"context"
	"errors"

	"gorm.io/gorm"
)

type UsersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {

	const op = "repository.postgres.GetUserByEmail"

	var user models.User
	if err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).
		Error; err != nil {

		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return &user, nil
}

func (r *UsersRepository) GetAllUsers(
	ctx context.Context,
	page, size int,
) ([]models.User, int64, error) {

	const op = "repository.postgres.GetAllUsers"

	offset := page * size

	var (
		users []models.User
		total int64
	)

	query := r.db.WithContext(ctx).Model(&models.User{})

	if err := query.Count(&total).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Count error: %v",
			op, observability.GetTraceID(ctx), err)
		return nil, 0, TranslateGormError(err)
	}

	if err := query.
		Limit(size).
		Offset(offset).
		Order("created_at ASC").
		Find(&users).
		Error; err != nil {

		logger.Error.Printf("[%s] TraceId=%s Find error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, 0, TranslateGormError(err)
	}

	return users, total, nil
}

func (r *UsersRepository) GetUserByID(ctx context.Context, userID string) (user *models.User, err error) {
	const op = "repository.postgres.GetUserByID"

	if err = r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}

		return nil, TranslateGormError(err)
	}

	return user, nil
}

func (r *UsersRepository) GetUsersByTeamID(ctx context.Context, teamID string) (users *[]models.User, err error) {
	const op = "repository.postgres.GetUsersByTeamID"

	if err = r.db.WithContext(ctx).Where("team_id = ?", teamID).Find(&users).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return users, nil
}

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user *models.User,
) (*models.User, error) {

	const op = "repository.postgres.CreateUser"

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return user, nil
}

func (r *UsersRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	const op = "repository.postgres.UpdateUser"

	if err := r.db.WithContext(ctx).Model(user).
		Where("id = ?", user.ID).
		Updates(user).
		First(user).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return user, nil
}

func (r *UsersRepository) UpdateUserStatus(ctx context.Context, userID string, isActive bool) error {
	const op = "repository.postgres.UpdateUserStatus"

	result := r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_active", isActive)

	if result.Error != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), result.Error)

		return TranslateGormError(result.Error)
	}

	return nil
}

func (r *UsersRepository) IsActive(ctx context.Context, userID string) (bool, error) {
	const op = "repository.postgres.IsActive"

	var user models.User
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return false, TranslateGormError(err)
	}

	return user.IsActive, nil
}
