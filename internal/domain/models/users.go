package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	// Admin - управляет пользователями, настраивает правила ревью
	// Может всё: создавать/редактировать пользователей, назначать роли,
	// настраивать кто и сколько должен одобрять эксперименты
	RoleAdmin Role = "admin"

	// Experimenter - создаёт и ведёт эксперименты
	// Отправляет эксперименты на ревью, получает решения аппруверов
	// Основной "пользователь" платформы
	RoleExperimenter Role = "experimenter"

	// Approver - проверяет и одобряет эксперименты перед запуском
	// Может одобрять/отклонять эксперименты, оставлять комментарии
	// Защищает от случайного запуска опасных экспериментов
	RoleApprover Role = "approver"

	// Viewer - только просмотр
	// Может смотреть эксперименты, отчеты, но ничего не менять
	// Ни создавать, ни ревьюить, ни администрировать
	RoleViewer Role = "viewer"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string         `gorm:"type:varchar(254);uniqueIndex;not null"`
	FullName  string         `gorm:"type:varchar(200);not null"`
	Password  string         `gorm:"type:varchar(255);not null"`
	Region    *string        `gorm:"type:varchar(32)"`
	Age       *int           `gorm:"type:integer"`
	Role      Role           `gorm:"type:varchar(20);"`
	TeamID    *uuid.UUID     `gorm:"index"`
	IsActive  bool           `gorm:"default:true"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
