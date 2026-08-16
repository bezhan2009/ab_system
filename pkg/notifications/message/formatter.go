package message

import (
	"ab_system/pkg/notifications/dto"
	"fmt"
	"strings"
	"time"
)

type EventType string

const (
	EventExperimentStarted   EventType = "experiment_started"
	EventExperimentPaused    EventType = "experiment_paused"
	EventExperimentCompleted EventType = "experiment_completed"
	EventGuardrailTriggered  EventType = "guardrail_triggered"
	EventExperimentReview    EventType = "experiment_review"
	EventExperimentApproved  EventType = "experiment_approved"
	EventExperimentRejected  EventType = "experiment_rejected"
)

var humanNames = map[EventType]string{
	EventExperimentStarted:   "Эксперимент запущен",
	EventExperimentPaused:    "Эксперимент на паузе",
	EventExperimentCompleted: "Эксперимент завершён",
	EventGuardrailTriggered:  "Сработал guardrail",
	EventExperimentReview:    "Эксперимент на ревью",
	EventExperimentApproved:  "Эксперимент одобрен",
	EventExperimentRejected:  "Эксперимент отклонён",
}

func escapeMarkdownV2(s string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`_`, `\_`,
		`*`, `\*`,
		`[`, `\[`,
		`]`, `\]`,
		`(`, `\(`,
		`)`, `\)`,
		`~`, `\~`,
		"`", "\\`",
		`>`, `\>`,
		`#`, `\#`,
		`+`, `\+`,
		`-`, `\-`,
		`=`, `\=`,
		`|`, `\|`,
		`{`, `\{`,
		`}`, `\}`,
		`.`, `\.`,
		`!`, `\!`,
	)

	return replacer.Replace(s)
}

func BuildTelegramMessage(req *dto.NotifyRequest) string {
	eventType := EventType(req.EventType)

	title, ok := humanNames[eventType]
	if !ok {
		title = string(eventType)
	}

	expID := req.ExperimentID
	if expID == "" {
		expID = "—"
	}
	expName := req.Experiment
	if expName == "" {
		expName = "—"
	}
	flagKey := req.FlagKey
	if flagKey == "" {
		flagKey = "—"
	}
	actor := req.UserID
	if actor == "" {
		actor = "—"
	}
	comment := req.Message
	url := req.ReportURL

	tsStr := time.Now().UTC().Format("02.01.2006 15:04 UTC")
	if req.Timestamp != "" {
		if dt, err := time.Parse(time.RFC3339, strings.Replace(req.Timestamp, "Z", "+00:00", 1)); err == nil {
			tsStr = dt.Format("02.01.2006 15:04 UTC")
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<b>%s</b>\n\n", title))
	sb.WriteString(fmt.Sprintf("Эксперимент: <code>%s</code> (<code>%s</code>)\n", expName, expID))
	sb.WriteString(fmt.Sprintf("Флаг: <code>%s</code>\n", flagKey))
	sb.WriteString(fmt.Sprintf("Кто: %s\n", actor))
	sb.WriteString(fmt.Sprintf("Время: <code>%s</code>\n", tsStr))

	details := make(map[string]string)
	switch eventType {
	case EventExperimentCompleted:
		if req.Status != "" {
			details["Результат"] = req.Status
		}
	case EventGuardrailTriggered:
		if req.Metric != "" {
			details["Метрика"] = req.Metric
		}
		if req.Threshold != "" {
			details["Порог"] = req.Threshold
		}
		if req.Value != "" {
			details["Значение"] = req.Value
		}
		if req.Status != "" {
			details["Действие"] = req.Status
		}
	}

	if len(details) > 0 {
		sb.WriteString("\n<b>Детали:</b>\n")
		for k, v := range details {
			sb.WriteString(fmt.Sprintf(" • %s: <code>%s</code>\n", k, v))
		}
	}

	if comment != "" {
		sb.WriteString(fmt.Sprintf("\n<i>%s</i>\n", comment))
	}

	if url != "" {
		sb.WriteString(fmt.Sprintf("\n<a href=\"%s\">Ссылка для просмотра отчёта: <code>%s</code></a>", url, url))
	}

	return sb.String()
}

func BuildSlackMessage(req *dto.NotifyRequest) string {
	eventType := EventType(req.EventType)

	title, ok := humanNames[eventType]
	if !ok {
		title = string(eventType)
	}

	expID := req.ExperimentID
	if expID == "" {
		expID = "—"
	}

	expName := req.Experiment
	if expName == "" {
		expName = "—"
	}

	flagKey := req.FlagKey
	if flagKey == "" {
		flagKey = "—"
	}

	actor := req.UserID
	if actor == "" {
		actor = "—"
	}

	comment := req.Message
	url := req.ReportURL

	tsStr := time.Now().UTC().Format("02.01.2006 15:04 UTC")
	if req.Timestamp != "" {
		if dt, err := time.Parse(time.RFC3339, strings.Replace(req.Timestamp, "Z", "+00:00", 1)); err == nil {
			tsStr = dt.Format("02.01.2006 15:04 UTC")
		}
	}

	lines := []string{
		fmt.Sprintf("*%s*", title),
		"",
		fmt.Sprintf("Эксперимент: %s (%s)", expName, expID),
		fmt.Sprintf("Флаг: %s", flagKey),
		fmt.Sprintf("Кто: %s", actor),
		fmt.Sprintf("Время: %s", tsStr),
	}

	details := make(map[string]string)

	switch eventType {
	case EventExperimentCompleted:
		if req.Status != "" {
			details["Результат"] = req.Status
		}
	case EventGuardrailTriggered:
		if req.Metric != "" {
			details["Метрика"] = req.Metric
		}
		if req.Threshold != "" {
			details["Порог"] = req.Threshold
		}
		if req.Value != "" {
			details["Значение"] = req.Value
		}
		if req.Status != "" {
			details["Действие"] = req.Status
		}
	}

	if len(details) > 0 {
		lines = append(lines, "", "*Детали:*")
		for k, v := range details {
			lines = append(lines, fmt.Sprintf("• %s: %s", k, v))
		}
	}

	if comment != "" {
		lines = append(lines, "", fmt.Sprintf("_%s_", comment))
	}

	if url != "" {
		lines = append(lines, "", fmt.Sprintf("Ссылка для просмотра отчёта: %s", url))
	}

	return strings.Join(lines, "\n")
}
