package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ExperimentStatus string

const (
	// Черновик - эксперимент создан, но еще не отправлен на ревью
	// Можно редактировать любые поля
	StatusDraft ExperimentStatus = "draft"

	// На ревью - эксперимент отправлен на проверку аппруверам
	// Ждет одобрения/отклонения
	StatusInReview ExperimentStatus = "in_review"

	// Одобрен - прошел ревью, набрал нужное количество одобрений
	// Готов к запуску, но еще не запущен
	StatusApproved ExperimentStatus = "approved"

	// Запущен - активен, участвует в выдаче значений флага
	// Пользователи получают варианты из этого эксперимента
	StatusRunning ExperimentStatus = "running"

	// На паузе - временно остановлен
	// Не участвует в выдаче, но конфигурация сохраняется
	// Можно возобновить (вернуть в running)
	StatusPaused ExperimentStatus = "paused"

	// Завершён - эксперимент закончен, решение принято
	// Дальше только архивация
	StatusCompleted ExperimentStatus = "completed"

	// В архиве - для истории, больше никогда не запустится
	StatusArchived ExperimentStatus = "archived"

	// Отклонён - не прошел ревью
	// Можно доработать и отправить заново (вернуть в draft)
	StatusRejected ExperimentStatus = "rejected"
)

type Experiment struct {
	ID      uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title   string           `gorm:"size:255;not null"`
	FlagKey string           `gorm:"size:255;not null;index"`
	Status  ExperimentStatus `gorm:"size:50;not null;default:'draft'"`

	Version int `gorm:"not null;default:1"`

	Conclusion      string `gorm:"size:255;not null;default:''"`
	Comment         string `gorm:"size:1000;not null;default:''"`
	WinnerVariantID string `gorm:"size:255;not null;default:''"`

	TargetingDsl   string `gorm:"type:text"`
	TrafficPercent int    `gorm:"not null"`

	GuardrailTriggered bool `gorm:"default:false"`
	RolledBackAt       *time.Time

	OwnerID string `gorm:"size:255;not null"`

	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
	StartedAt   *time.Time
	CompletedAt *time.Time

	Variants             []Variant            `gorm:"foreignKey:ExperimentID"`
	RampUp               ExperimentRampUp     `gorm:"foreignKey:ExperimentID"`
	NotificationSettings NotificationSettings `gorm:"foreignKey:ExperimentID"`
}

type ExperimentVersion struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ExperimentID uuid.UUID `gorm:"type:uuid;not null;index"`
	Version      int       `gorm:"not null"`

	Snapshot datatypes.JSON `gorm:"type:json;not null"`

	ChangedBy string `gorm:"size:255;not null"`

	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
