package configs

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

type Configs struct {
	LogParams          LogParams          `json:"log_params"`
	AppParams          AppParams          `json:"app_params"`
	ServerParams       ServerParams       `json:"server_params"`
	GuardrailParams    GuardrailParams    `json:"guardrail_params"`
	AutopilotParams    AutopilotParams    `json:"autopilot_params"`
	Clients            Clients            `json:"clients"`
	NotificationParams NotificationParams `json:"notification_params"`
}

type LogParams struct {
	LogDirectory     string `json:"log_directory"`
	LogInfo          string `json:"log_info"`
	LogError         string `json:"log_error"`
	LogWarn          string `json:"log_warn"`
	LogDebug         string `json:"log_debug"`
	MaxSizeMegabytes int    `json:"max_size_megabytes"`
	MaxBackups       int    `json:"max_backups"`
	MaxAge           int    `json:"max_age"`
	Compress         bool   `json:"compress"`
	LocalTime        bool   `json:"local_time"`
}

type AppParams struct {
	GinMode    string `json:"gin_mode"`
	PortRun    string `json:"port_run"`
	ServerURL  string `json:"server_url"`
	ServerName string `json:"server_name"`
}

type ServerParams struct {
	Addr                string `json:"addr"`
	MaxHeaderMBs        int    `json:"max_header_mbs"`
	ReadTimeoutSeconds  int    `json:"read_timeout_seconds"`
	WriteTimeoutSeconds int    `json:"write_timeout_seconds"`
}

type GuardrailParams struct {
	CheckerTicketMinutes int `json:"checker_ticket_minutes"`
}

type AutopilotParams struct {
	CheckerTicketMinutes int `json:"checker_ticket_minutes"`
}

type TelegramNotifications struct {
	Address string `json:"address"`
}

type NotificationParams struct {
	TtlNotificationsMinutes int `json:"ttl_notifications_minutes"`
}

type Clients struct {
	Telegram TelegramNotifications `json:"telegram_notifications"`
}

func ReadSettings() (Configs, error) {
	var AppSettings Configs

	configFile, err := os.Open(os.Getenv("CONFIG_PATH"))
	if err != nil {
		return Configs{}, errors.New(fmt.Sprintf("Couldn't open config file. Error is: %s", err.Error()))
	}

	defer func(configFile *os.File) {
		err = configFile.Close()
		if err != nil {
			log.Fatal("Couldn't close config file. Error is: ", err.Error())
		}
	}(configFile)

	if err = json.NewDecoder(configFile).Decode(&AppSettings); err != nil {
		return Configs{}, errors.New(fmt.Sprintf("Couldn't decode settings json file. Error is: %s", err.Error()))
	}
	return AppSettings, nil
}
