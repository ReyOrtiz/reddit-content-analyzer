package main

import (
	"fmt"

	"github.com/ReyOrtiz/reddit-content-analyzer/internal/api"
	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/config"
	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/logger"
	"go.uber.org/zap"
)

func main() {
	logger := logger.GetLogger()
	defer logger.Sync()
	cfg := config.GetConfig()
	port := cfg.GetString("api.port")

	logger.Info(
		"Starting Reddit Content Analyzer API",
		zap.String("version", "1.0.0"),
		zap.String("port", fmt.Sprintf(":%s", port)),
	)

	api.StartServer()
}
