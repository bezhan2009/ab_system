package jobs

import (
	"ab_system/internal/clients/telegram_notifications"
	"ab_system/internal/configs"
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	"ab_system/internal/domain/service"
	"ab_system/internal/lib/metrics"
	"ab_system/internal/notifications/slack"
	telegramv1 "ab_system/internal/telegram_bot/gen/telegram_notifications/v1"
	"ab_system/pkg/logger"
	pkgDto "ab_system/pkg/notifications/dto"
	"context"
	"fmt"
	"time"
)

type GuardrailChecker struct {
	appParams configs.AppParams

	experimentReader repository.ExperimentReader
	expMetricReader  repository.ExperimentMetricReader
	metricReader     repository.MetricReader
	eventReader      repository.EventReader
	decisionReader   repository.DecisionReader
	experimentWriter repository.ExperimentWriter
	triggerWriter    repository.GuardrailTriggerWriter
	metricLib        *metrics.MetricLib
	userReader       repository.UserReader

	notificationSettingsReader repository.NotificationSettingsReader

	telegramClient *telegram_notifications.NotifyClient
	slackNotifier  *slack.Notifier
	interval       time.Duration
}

func NewGuardrailChecker(
	appParams configs.AppParams,

	experimentReader repository.ExperimentReader,
	expMetricReader repository.ExperimentMetricReader,
	metricReader repository.MetricReader,
	eventReader repository.EventReader,
	decisionReader repository.DecisionReader,
	experimentWriter repository.ExperimentWriter,
	triggerWriter repository.GuardrailTriggerWriter,
	metricLib *metrics.MetricLib,
	userReader repository.UserReader,

	notificationSettingsReader repository.NotificationSettingsReader,

	telegramClient *telegram_notifications.NotifyClient,
	slackNotifier *slack.Notifier,
	interval time.Duration,
) *GuardrailChecker {
	return &GuardrailChecker{
		appParams: appParams,

		experimentReader: experimentReader,
		expMetricReader:  expMetricReader,
		metricReader:     metricReader,
		eventReader:      eventReader,
		decisionReader:   decisionReader,
		experimentWriter: experimentWriter,
		triggerWriter:    triggerWriter,
		metricLib:        metricLib,
		userReader:       userReader,

		notificationSettingsReader: notificationSettingsReader,

		telegramClient: telegramClient,
		slackNotifier:  slackNotifier,
		interval:       interval,
	}
}

func (c *GuardrailChecker) Start(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info.Println("Quitting Guardrail Checker")

				ticker.Stop()
				return
			case <-ticker.C:
				logger.Info.Println("Checking Guardrail Checker")

				c.checkAllExperiments(ctx)
			}
		}
	}()
}

func (c *GuardrailChecker) checkAllExperiments(ctx context.Context) {
	experiments, err := c.experimentReader.GetExperimentsByStatus(ctx, string(models.StatusRunning))
	if err != nil {
		logger.Info.Printf("GuardrailChecker: failed to get running experiments: %v", err)
		logger.Error.Printf("[service.GuardrailChecker.checkAllExperiments] Error getting running experiments: %v", err)
		return
	}

	for _, exp := range *experiments {
		c.checkExperiment(ctx, &exp)
	}
}

func (c *GuardrailChecker) checkExperiment(ctx context.Context, exp *models.Experiment) {
	guardrails, err := c.expMetricReader.GetGuardrailsForExperiment(ctx, exp.ID.String())
	if err != nil {
		logger.Info.Printf("GuardrailChecker: failed to get guardrails for experiment %s: %v", exp.ID, err)
		return
	}

	for _, gr := range guardrails {
		to := time.Now().UTC()
		from := to.Add(-time.Duration(gr.WindowMin) * time.Minute)

		decisions, err := c.decisionReader.GetDecisionsByExperimentAndTime(ctx, exp.ID.String(), from, to)
		if err != nil {
			logger.Info.Printf("GuardrailChecker: failed to get decisions: %v", err)
			continue
		}

		decIDs := make([]string, len(decisions))
		for i, d := range decisions {
			decIDs[i] = d.ID.String()
		}

		metric, err := c.metricReader.GetMetricByID(ctx, gr.MetricID.String())
		if err != nil {
			continue
		}

		value, err := c.metricLib.CalculateMetric(ctx, metric, decIDs, from, to, false)
		if err != nil {
			continue
		}

		valuePercent := value * 100

		if c.isThresholdExceeded(valuePercent, gr.Threshold, gr.Operator) {
			c.executeAction(ctx, exp, &gr, metric, value)
		}
	}
}

func (c *GuardrailChecker) isThresholdExceeded(value, threshold float64, operator string) bool {
	switch operator {
	case ">":
		return value > threshold
	case ">=":
		return value >= threshold
	case "<":
		return value < threshold
	case "<=":
		return value <= threshold
	default:
		return false
	}
}

func (c *GuardrailChecker) executeAction(ctx context.Context, exp *models.Experiment, gr *models.ExperimentMetric, metric *models.Metric, actualValue float64) {
	ownerEmail := ""

	owner, err := c.userReader.GetUserByID(ctx, exp.OwnerID)
	if err != nil {
		ownerEmail = exp.OwnerID
	} else {
		ownerEmail = fmt.Sprintf("%s (%s)", owner.Email, exp.OwnerID)
	}

	trigger := &models.GuardrailTrigger{
		ExperimentID: exp.ID,
		MetricID:     gr.MetricID,
		Threshold:    gr.Threshold,
		Operator:     gr.Operator,
		WindowMin:    gr.WindowMin,
		ActualValue:  actualValue * 100,
		Action:       gr.Action,
		TriggeredAt:  time.Now(),
	}
	_ = c.triggerWriter.CreateTrigger(ctx, trigger)

	switch gr.Action {
	case "pause":
		exp.Status = models.StatusPaused
		_, _ = c.experimentWriter.UpdateExperiment(ctx, exp)
	case "rollback":
		exp.GuardrailTriggered = true
		now := time.Now()
		exp.RolledBackAt = &now
		_, _ = c.experimentWriter.UpdateExperiment(ctx, exp)
	}

	notifyReq := &telegramv1.NotifyRequest{
		EventType:    "guardrail_triggered",
		ExperimentId: exp.ID.String(),
		Experiment:   exp.Title,
		FlagKey:      exp.FlagKey,
		UserId:       ownerEmail,
		Metric:       metric.Title,
		Threshold:    fmt.Sprintf("%.2f%%", gr.Threshold),
		Value:        fmt.Sprintf("%.2f%%", actualValue*100),
		Status:       gr.Action,
		Timestamp:    time.Now().Format(time.RFC3339),
		ReportUrl: fmt.Sprintf("http://%s:%s/api/v1/reports/%s",
			c.appParams.ServerURL, c.appParams.PortRun, exp.ID.String()),
	}

	go func() {
		bgCtx := context.Background()
		_ = service.SendTelegramNotification(bgCtx, c.telegramClient.Notify, notifyReq)
	}()

	go func() {
		bgCtx := context.Background()

		webhooks, err := c.notificationSettingsReader.GetSlackWebhooksForExpNotification(bgCtx, exp.ID.String())
		if err != nil {
			logger.Error.Printf("Failed to get slack webhooks for experiment %s: %v", exp.ID, err)
			return
		}

		slackReq := pkgDto.NotifyRequest{
			EventType:    notifyReq.EventType,
			ExperimentID: notifyReq.ExperimentId,
			Experiment:   notifyReq.Experiment,
			FlagKey:      notifyReq.FlagKey,
			UserID:       notifyReq.UserId,
			Status:       notifyReq.Status,
			Metric:       notifyReq.Metric,
			Threshold:    notifyReq.Threshold,
			Value:        notifyReq.Value,
			Message:      notifyReq.Message,
			ReportURL:    notifyReq.ReportUrl,
			Timestamp:    notifyReq.Timestamp,
		}

		if len(webhooks) > 0 {
			if err := c.slackNotifier.Send(bgCtx, slackReq, webhooks); err != nil {
				logger.Error.Printf("Failed to send slack notification for experiment %s: %v", exp.ID, err)
			}
		}
	}()
}
