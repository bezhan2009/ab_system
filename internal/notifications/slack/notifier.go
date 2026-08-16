package slack

import (
	"ab_system/internal/domain/redis"
	"ab_system/pkg/logger"
	notifydto "ab_system/pkg/notifications/dto"
	"ab_system/pkg/notifications/message"
	"ab_system/pkg/utils"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Notifier struct {
	httpClient *http.Client
	redisCache redis.RedisCacheRepository
	ttl        time.Duration
}

func NewNotifier(redisCache redis.RedisCacheRepository, ttl time.Duration) *Notifier {
	return &Notifier{
		httpClient: &http.Client{Timeout: 60 * time.Second},
		redisCache: redisCache,
		ttl:        ttl,
	}
}

func (n *Notifier) Send(ctx context.Context, req notifydto.NotifyRequest, webhooks []string) (err error) {
	const op = "slack.Send"

	if len(webhooks) == 0 {
		return nil
	}

	key := utils.GenerateNotificationKey("slack", &req)

	success, err := n.redisCache.SetNX(ctx, key, "1", n.ttl)
	if err != nil {
		logger.Error.Printf("[%s] Cache error: %v", op, err)
	} else if !success {
		logger.Info.Printf("[%s] Duplicate Slack notification skipped, key: %s", op, key)

		return nil
	}

	text := message.BuildSlackMessage(&req)

	payload := map[string]interface{}{
		"text": text,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error.Printf("[%s] Error marshaling Slack payload: %v", op, err)
		return err
	}

	for _, webhook := range webhooks {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhook, bytes.NewReader(jsonPayload))
		if err != nil {
			logger.Error.Printf("[%s] Error creating request for webhook %s: %v", op, webhook, err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := n.httpClient.Do(req)
		if err != nil {
			logger.Error.Printf("[%s] Error sending to Slack webhook %s: %v", op, webhook, err)
			continue
		}

		resp.Body.Close()

		if resp.StatusCode >= 300 {
			logger.Error.Printf("[%s] Slack webhook %s returned status %d", op, webhook, resp.StatusCode)
		} else {
			logger.Info.Printf("[%s] Sent Slack notification to %s", op, webhook)
		}
	}

	return nil
}
