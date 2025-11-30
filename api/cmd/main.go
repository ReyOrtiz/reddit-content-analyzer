package main

import (
	"fmt"

	"github.com/ReyOrtiz/reddit-content-analyzer/internal/api"
	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/config"
	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/logger"
	"go.uber.org/zap"
)

// @title           Reddit Content Analyzer API
// @version         1.0
// @description     API for analyzing and searching relevant Reddit content based on topics
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /v1

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
