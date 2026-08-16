package seeds

import (
	"ab_system/internal/domain/models"
	"ab_system/pkg/logger"
	"errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func SeedEventTypes(db *gorm.DB) error {
	eventTypes := []models.EventType{
		{
			Title:              "impression",
			Description:        "Показ интерфейса пользователю",
			Schema:             datatypes.JSON(`{"screen": "string", "element": "string"}`),
			RequiresDecisionID: true,
			RequiresUserID:     true,
			RequiresExposure:   false,
		},
		{
			Title:              "click",
			Description:        "Клик по элементу интерфейса",
			Schema:             datatypes.JSON(`{"element": "string", "screen": "string"}`),
			RequiresDecisionID: true,
			RequiresUserID:     true,
			RequiresExposure:   true,
		},
		{
			Title:              "purchase",
			Description:        "Совершение покупки",
			Schema:             datatypes.JSON(`{"amount": "number", "currency": "string", "items_count": "number"}`),
			RequiresDecisionID: true,
			RequiresUserID:     true,
			RequiresExposure:   true,
		},
		{
			Title:              "error",
			Description:        "Ошибка в приложении",
			Schema:             datatypes.JSON(`{"error_code": "string", "error_message": "string"}`),
			RequiresDecisionID: true,
			RequiresUserID:     true,
			RequiresExposure:   true,
		},
		{
			Title:              "api_response",
			Description:        "Время ответа API",
			Schema:             datatypes.JSON(`{"duration_ms": "number", "endpoint": "string"}`),
			RequiresDecisionID: true,
			RequiresUserID:     true,
			RequiresExposure:   true,
		},
		{
			Title:              "page_view",
			Description:        "Просмотр страницы",
			Schema:             datatypes.JSON(`{"page": "string", "referrer": "string"}`),
			RequiresDecisionID: true,
			RequiresUserID:     true,
			RequiresExposure:   false,
		},
	}

	for _, et := range eventTypes {
		var existing models.EventType
		err := db.Where("title = ?", et.Title).First(&existing).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&et).Error; err != nil {
					logger.Error.Printf("[seeds.SeedEventTypes] Error creating event type %s: %v", et.Title, err)
					return err
				}
				logger.Info.Printf("[seeds.SeedEventTypes] Created event type: %s", et.Title)
			} else {
				logger.Error.Printf("[seeds.SeedEventTypes] Error checking event type %s: %v", et.Title, err)
				return err
			}
		}
	}

	return nil
}
