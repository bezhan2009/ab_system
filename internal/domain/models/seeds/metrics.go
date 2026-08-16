package seeds

import (
	"ab_system/internal/domain/models"
	"ab_system/pkg/logger"
	"errors"

	"gorm.io/gorm"
)

func SeedMetrics(db *gorm.DB) error {
	metrics := []models.Metric{
		{
			Title:            "impression_count",
			Description:      "Количество показов",
			Type:             "counter",
			CounterEventType: "impression",
			RequiresExposure: false,
			Unit:             "count",
		},
		{
			Title:            "click_count",
			Description:      "Количество кликов",
			Type:             "counter",
			CounterEventType: "click",
			RequiresExposure: true,
			Unit:             "count",
		},
		{
			Title:            "purchase_count",
			Description:      "Количество покупок",
			Type:             "counter",
			CounterEventType: "purchase",
			RequiresExposure: true,
			Unit:             "count",
		},
		{
			Title:            "error_count",
			Description:      "Количество ошибок",
			Type:             "counter",
			CounterEventType: "error",
			RequiresExposure: true,
			Unit:             "count",
		},

		{
			Title:                "ctr",
			Description:          "CTR (клики/показы)",
			Type:                 "ratio",
			NumeratorEventType:   "click",
			DenominatorEventType: "impression",
			RequiresExposure:     true,
			Unit:                 "%",
		},
		{
			Title:                "conversion_rate",
			Description:          "Конверсия в покупку",
			Type:                 "ratio",
			NumeratorEventType:   "purchase",
			DenominatorEventType: "impression",
			RequiresExposure:     true,
			Unit:                 "%",
		},
		{
			Title:                "error_rate",
			Description:          "Доля ошибок",
			Type:                 "ratio",
			NumeratorEventType:   "error",
			DenominatorEventType: "impression",
			RequiresExposure:     true,
			Unit:                 "%",
		},

		{
			Title:              "avg_latency",
			Description:        "Среднее время ответа",
			Type:               "histogram",
			HistogramEventType: "api_response",
			HistogramField:     "duration_ms",
			RequiresExposure:   true,
			Unit:               "ms",
		},
		{
			Title:              "p95_latency",
			Description:        "P95 время ответа",
			Type:               "histogram",
			HistogramEventType: "api_response",
			HistogramField:     "duration_ms",
			RequiresExposure:   true,
			Unit:               "ms",
		},
		{
			Title:              "avg_purchase_amount",
			Description:        "Средняя сумма покупки",
			Type:               "histogram",
			HistogramEventType: "purchase",
			HistogramField:     "amount",
			RequiresExposure:   true,
			Unit:               "RUB",
		},
	}

	for _, m := range metrics {
		var existing models.Metric
		err := db.Where("title = ?", m.Title).First(&existing).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err = db.Create(&m).Error; err != nil {
					logger.Error.Printf("[seeds.SeedMetrics] Error creating metric %s: %v", m.Title, err)
					return err
				}
				logger.Info.Printf("[seeds.SeedMetrics] Created metric: %s", m.Title)
			} else {
				logger.Error.Printf("[seeds.SeedMetrics] Error checking metric %s: %v", m.Title, err)
				return err
			}
		}
	}

	return nil
}
