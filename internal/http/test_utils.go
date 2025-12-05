package http

import (
	"log/slog"
	"os"
	"task-manager/internal/config"
	"task-manager/internal/service"
	"task-manager/internal/utils"
	"task-manager/pkg/logger"
)

func SetupHandler(taskService *service.MockTaskService) *Handler {
	consoleHandler := slog.NewTextHandler(os.Stdout, nil)

	slogLogger := slog.New(consoleHandler)

	myLogger := &logger.StandardLogger{
		Logger: slogLogger,
	}

	return CreateHandler(myLogger, config.Config{}, taskService, utils.InitGlobalTaskMetrics())
}
