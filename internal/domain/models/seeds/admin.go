package seeds

import (
	"ab_system/internal/domain/models"
	"ab_system/pkg/logger"
	"ab_system/pkg/utils"
	"errors"
	"os"

	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) (err error) {
	fullName := os.Getenv("ADMIN_FULLNAME")
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")

	var admin models.User

	admin.Email = email

	admin.IsActive = true

	admin.Password, err = utils.HashPassword(password)
	if err != nil {
		return err
	}

	admin.FullName = fullName
	admin.Role = models.RoleAdmin

	var existingAdmin models.User
	if err = db.Where("email = ?", email).First(&existingAdmin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = db.Create(&admin).Error; err != nil {
				logger.Error.Printf("[seeds.SeedRoles] Error creating role %s: %v", admin.Email, err)

				return err
			}
		} else {
			logger.Error.Printf("[seeds.SeedRoles] Error checking role %s: %v", admin.Email, err)

			return err
		}
	} else {
		if !utils.CheckPassword(admin.Password, existingAdmin.Password) {
			err = db.Where("id = ?", existingAdmin.ID).Updates(&admin).Error
			if err != nil {
				logger.Error.Printf("[seeds.SeedRoles] Error updating role %s: %v", admin.Email, err)
				return err
			}
		}
	}

	return nil
}
