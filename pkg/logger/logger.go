package logger

import (
	"ab_system/internal/configs"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Info  *log.Logger
	Error *log.Logger
	Warn  *log.Logger
	Debug *log.Logger
)

func InitLogger(logParams configs.LogParams) error {
	if _, err := os.Stat(logParams.LogDirectory); os.IsNotExist(err) {
		err = os.Mkdir(logParams.LogDirectory, 0755)
		if err != nil {
			return err
		}
	}

	lumberLogInfo := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", logParams.LogDirectory, logParams.LogInfo),
		MaxSize:    logParams.MaxSizeMegabytes,
		MaxBackups: logParams.MaxBackups,
		MaxAge:     logParams.MaxAge,
		Compress:   logParams.Compress,
		LocalTime:  logParams.LocalTime,
	}

	lumberLogError := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", logParams.LogDirectory, logParams.LogError),
		MaxSize:    logParams.MaxSizeMegabytes,
		MaxBackups: logParams.MaxBackups,
		MaxAge:     logParams.MaxAge,
		Compress:   logParams.Compress,
		LocalTime:  logParams.LocalTime,
	}

	lumberLogWarn := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", logParams.LogDirectory, logParams.LogWarn),
		MaxSize:    logParams.MaxSizeMegabytes,
		MaxBackups: logParams.MaxBackups,
		MaxAge:     logParams.MaxAge,
		Compress:   logParams.Compress,
		LocalTime:  logParams.LocalTime,
	}

	lumberLogDebug := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", logParams.LogDirectory, logParams.LogDebug),
		MaxSize:    logParams.MaxSizeMegabytes,
		MaxBackups: logParams.MaxBackups,
		MaxAge:     logParams.MaxAge,
		Compress:   logParams.Compress,
		LocalTime:  logParams.LocalTime,
	}

	Info = log.New(gin.DefaultWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(lumberLogError, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(lumberLogWarn, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(lumberLogDebug, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	gin.DefaultWriter = io.MultiWriter(os.Stdout, lumberLogInfo)

	return nil
}
