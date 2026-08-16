package dto

import (
	"ab_system/internal/domain/models"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Variant struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Value       string `json:"value"`
	Weight      int    `json:"weight"`
	IsControl   bool   `json:"is_control"`
	Description string `json:"description,omitempty"`
}

type ExperimentRampUp struct {
	RampEnabled         bool      `json:"ramp_enabled"`
	RampSteps           []int     `json:"ramp_steps"`
	RampCurrentStep     int       `json:"ramp_current_step"`
	RampLastIncrease    time.Time `json:"ramp_last_increase"`
	RampIntervalMinutes int       `json:"ramp_interval_minutes"`
}

type Experiment struct {
	ID string `json:"id"`

	RolledBackAt *time.Time `json:"rolled_back_at"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	StartedAt    *time.Time `json:"startedAt,omitempty"`
	PausedAt     *time.Time `json:"pausedAt,omitempty"`
	CompletedAt  *time.Time `json:"completedAt,omitempty"`

	Title       string `json:"title"`
	Description string `json:"description"`
	FlagKey     string `json:"flag_key"`
	Status      string `json:"status"`

	Conclusion      string `json:"conclusion"`
	Comment         string `json:"comment"`
	WinnerVariantID string `json:"winner_variant_id"`

	TrafficPercent int    `json:"traffic_percent"`
	TargetingDsl   string `json:"targeting_dsl"`

	GuardrailTriggered bool `json:"guardrail_triggered"`

	Version int `json:"version"`

	OwnerID string `json:"owner_id"`

	Variants []*Variant `json:"variants"`

	RampEnabled         bool      `json:"ramp_enabled,omitempty"`
	RampSteps           []int     `json:"ramp_steps,omitempty"`
	RampCurrentStep     int       `json:"ramp_current_step,omitempty"`
	RampLastIncrease    time.Time `json:"ramp_last_increase,omitempty"`
	RampIntervalMinutes int       `json:"ramp_interval_minutes,omitempty"`

	NotificationSettings NotificationSettings `json:"notification_settings"`
}

type ExperimentVersion struct {
	ID           string          `json:"id"`
	ExperimentID string          `json:"experiment_id"`
	Version      int             `json:"version"`
	Snapshot     json.RawMessage `json:"snapshot"`
	ChangedBy    string          `json:"changed_by"`
	ChangeReason string          `json:"change_reason,omitempty"`
	CreatedAt    time.Time       `json:"createdAt"`
}

type CompleteExperimentRequest struct {
	Conclusion      string `json:"conclusion"`
	Comment         string `json:"comment"`
	WinnerVariantID string `json:"winner_variant_id"`
}

func (d *Experiment) ToModel() (*models.Experiment, error) {
	var err error

	expId := uuid.UUID{}
	expIdStr := d.ID
	if expIdStr != "" {
		expId, err = uuid.Parse(expIdStr)
		if err != nil {
			return nil, err
		}
	}

	ownerID, err := uuid.Parse(d.OwnerID)
	if err != nil {
		return nil, err
	}

	variants := make([]models.Variant, 0, len(d.Variants))
	for _, v := range d.Variants {
		variantID := uuid.UUID{}
		if v.ID != "" {
			variantID, err = uuid.Parse(v.ID)
			if err != nil {
				return nil, err
			}
		} else {
			variantID = uuid.New()
		}

		variants = append(variants, models.Variant{
			ID:          variantID,
			Title:       v.Title,
			Value:       v.Value,
			Weight:      v.Weight,
			IsControl:   v.IsControl,
			Description: v.Description,
		})
	}

	model := &models.Experiment{
		ID:                 expId,
		Conclusion:         d.Conclusion,
		Comment:            d.Comment,
		RolledBackAt:       d.RolledBackAt,
		CreatedAt:          d.CreatedAt,
		UpdatedAt:          d.UpdatedAt,
		StartedAt:          d.StartedAt,
		CompletedAt:        d.CompletedAt,
		Title:              d.Title,
		FlagKey:            d.FlagKey,
		TargetingDsl:       d.TargetingDsl,
		Status:             models.ExperimentStatus(d.Status),
		GuardrailTriggered: d.GuardrailTriggered,
		TrafficPercent:     d.TrafficPercent,
		OwnerID:            ownerID.String(),
		Version:            d.Version,
		Variants:           variants,
	}

	if d.RampEnabled {
		stepsJSON, _ := json.Marshal(d.RampSteps)
		model.RampUp = models.ExperimentRampUp{
			ExperimentID:        &model.ID,
			RampEnabled:         d.RampEnabled,
			RampSteps:           stepsJSON,
			RampCurrentStep:     d.RampCurrentStep,
			RampLastIncrease:    d.RampLastIncrease,
			RampIntervalMinutes: d.RampIntervalMinutes,
		}
	}

	notifId := uuid.Nil
	if d.NotificationSettings.ID != "" {
		notifId, err = uuid.Parse(d.NotificationSettings.ID)
		if err != nil {
			return nil, err
		}
	} else {
		notifId = uuid.New()
	}

	if d.NotificationSettings.ChatIds == nil {
		d.NotificationSettings.ChatIds = []int64{}
	}

	chatIdsJSON, _ := json.Marshal(d.NotificationSettings.ChatIds)
	slackWebhooksJSON, _ := json.Marshal(d.NotificationSettings.SlackWebhooks)

	model.NotificationSettings = models.NotificationSettings{
		ID:            notifId,
		ExperimentID:  &model.ID,
		ChatIds:       chatIdsJSON,
		SlackWebhooks: slackWebhooksJSON,
		CreatedAt:     d.NotificationSettings.CreatedAt,
		UpdatedAt:     d.NotificationSettings.UpdatedAt,
	}

	return model, nil
}

func (d *Experiment) ToDTO(exp *models.Experiment) *Experiment {
	variants := make([]*Variant, len(exp.Variants))
	for i, v := range exp.Variants {
		variants[i] = &Variant{
			ID:          v.ID.String(),
			Title:       v.Title,
			Value:       v.Value,
			Weight:      v.Weight,
			IsControl:   v.IsControl,
			Description: v.Description,
		}
	}

	dto := &Experiment{
		ID:                 exp.ID.String(),
		Title:              exp.Title,
		Conclusion:         exp.Conclusion,
		Comment:            exp.Comment,
		RolledBackAt:       exp.RolledBackAt,
		CreatedAt:          exp.CreatedAt,
		UpdatedAt:          exp.UpdatedAt,
		StartedAt:          exp.StartedAt,
		CompletedAt:        exp.CompletedAt,
		TargetingDsl:       exp.TargetingDsl,
		FlagKey:            exp.FlagKey,
		Status:             string(exp.Status),
		GuardrailTriggered: exp.GuardrailTriggered,
		TrafficPercent:     exp.TrafficPercent,
		OwnerID:            exp.OwnerID,
		Version:            exp.Version,
		Variants:           variants,
	}

	if exp.RampUp.ID != uuid.Nil {
		var steps []int
		json.Unmarshal(exp.RampUp.RampSteps, &steps)

		dto.RampEnabled = exp.RampUp.RampEnabled
		dto.RampSteps = steps
		dto.RampCurrentStep = exp.RampUp.RampCurrentStep
		dto.RampLastIncrease = exp.RampUp.RampLastIncrease
		dto.RampIntervalMinutes = exp.RampUp.RampIntervalMinutes
	}

	var chatIds []int64
	if exp.NotificationSettings.ChatIds != nil {
		if err := json.Unmarshal(exp.NotificationSettings.ChatIds, &chatIds); err != nil {
			chatIds = []int64{}
		}
	} else {
		chatIds = []int64{}
	}

	var slackWebhooks []string
	if exp.NotificationSettings.SlackWebhooks != nil {
		err := json.Unmarshal(exp.NotificationSettings.SlackWebhooks, &slackWebhooks)
		if err != nil {
			slackWebhooks = []string{}
		}
	} else {
		slackWebhooks = []string{}
	}

	dto.NotificationSettings = NotificationSettings{
		ID:            exp.NotificationSettings.ID.String(),
		ExperimentID:  exp.ID.String(),
		ChatIds:       chatIds,
		SlackWebhooks: slackWebhooks,
		CreatedAt:     exp.NotificationSettings.CreatedAt,
		UpdatedAt:     exp.NotificationSettings.UpdatedAt,
	}

	return dto
}

func (d *Experiment) ToDTOs(experiments []models.Experiment) []*Experiment {
	result := make([]*Experiment, len(experiments))
	for i := range experiments {
		result[i] = d.ToDTO(&experiments[i])
	}

	return result
}

func (d *ExperimentVersion) ToModel() (*models.ExperimentVersion, error) {
	id := uuid.UUID{}
	if d.ID != "" {
		var err error

		id, err = uuid.Parse(d.ID)
		if err != nil {
			return nil, err
		}
	}

	expID, err := uuid.Parse(d.ExperimentID)
	if err != nil {
		return nil, err
	}

	return &models.ExperimentVersion{
		ID:           id,
		ExperimentID: expID,
		Version:      d.Version,
		Snapshot:     datatypes.JSON(d.Snapshot),
		ChangedBy:    d.ChangedBy,
		CreatedAt:    d.CreatedAt,
	}, nil
}

func (d *ExperimentVersion) ToDTO(version *models.ExperimentVersion) *ExperimentVersion {
	return &ExperimentVersion{
		ID:           version.ID.String(),
		ExperimentID: version.ExperimentID.String(),
		Version:      version.Version,
		Snapshot:     json.RawMessage(version.Snapshot),
		ChangedBy:    version.ChangedBy,
		CreatedAt:    version.CreatedAt,
	}
}

func (d *ExperimentVersion) ToDTOs(versions []models.ExperimentVersion) []*ExperimentVersion {
	result := make([]*ExperimentVersion, len(versions))
	for i := range versions {
		result[i] = d.ToDTO(&versions[i])
	}

	return result
}
