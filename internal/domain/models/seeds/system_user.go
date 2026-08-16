package seeds

import (
	"ab_system/internal/domain/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedSystemUser(db *gorm.DB) (uuid.UUID, error) {
	systemID := uuid.MustParse("188d00c6-c033-41bc-9b71-0d7e9dcecf91")

	var user models.User
	err := db.Where("id = ?", systemID).First(&user).Error
	if err == nil {
		return systemID, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return uuid.Nil, err
	}

	systemUser := models.User{
		ID:       systemID,
		Email:    "system@internal",
		FullName: "System User",
		Password: "system",
		Role:     models.RoleAdmin,
		IsActive: true,
	}

	if err = db.Create(&systemUser).Error; err != nil {
		return uuid.Nil, err
	}

	return systemID, nil
}
