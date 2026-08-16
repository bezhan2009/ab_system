package postgres

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/models/seeds"
	"errors"

	"gorm.io/gorm"
)

func Migrate(dbConn *gorm.DB) error {
	if dbConn == nil {
		return errors.New("database connection is not initialized")
	}

	tables := []interface{}{
		&models.User{},
		&models.Team{},
		&models.FeatureFlag{},
		&models.Experiment{},
		&models.ExperimentVersion{},
		&models.ExperimentRampUp{},
		&models.Variant{},
		&models.Decision{},
		&models.ApproverGroup{},
		&models.ApproverGroupMember{},
		&models.Approval{},
		&models.EventType{},
		&models.Event{},
		&models.Metric{},
		&models.ExperimentMetric{},
		&models.GuardrailTrigger{},
		&models.NotificationSettings{},
	}

	for _, table := range tables {
		if err := dbConn.AutoMigrate(table); err != nil {
			return err
		}
	}

	//systemUserId, err := seeds.SeedSystemUser(dbConn)
	//if err != nil {
	//	return err
	//}

	if err := seeds.SeedAdmin(dbConn); err != nil {
		return err
	}

	if err := seeds.SeedEventTypes(dbConn); err != nil {
		return err
	}

	if err := seeds.SeedMetrics(dbConn); err != nil {
		return err
	}

	return nil
}
