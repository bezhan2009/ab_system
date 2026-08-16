package handlers

import (
	"ab_system/internal/clients/telegram_notifications"
	"ab_system/internal/http/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingHandler struct {
	telegramClient *telegram_notifications.NotifyClient
}

func NewPingHandler(telegramClient *telegram_notifications.NotifyClient) *PingHandler {
	return &PingHandler{
		telegramClient: telegramClient,
	}
}

func (h *PingHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *PingHandler) Ready(c *gin.Context) {
	err := h.telegramClient.Ready(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}
