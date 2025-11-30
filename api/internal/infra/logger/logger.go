package logger

import (
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ReyOrtiz/reddit-content-analyzer/internal/infra/config"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// GetLogger returns the singleton logger instance, initializing it on first call
func GetLogger() *zap.Logger {
	once.Do(func() {
		cfg := config.GetConfig()

		// Get log level from config
		levelStr := strings.ToLower(cfg.GetString("logging.level"))
		var level zapcore.Level
		switch levelStr {
		case "debug":
			level = zapcore.DebugLevel
		case "info":
			level = zapcore.InfoLevel
		case "warn":
			level = zapcore.WarnLevel
		case "error":
			level = zapcore.ErrorLevel
		default:
			level = zapcore.InfoLevel
		}

		// Get log format from config
		format := strings.ToLower(cfg.GetString("logging.format"))

		// Configure encoder based on format
		var encoderConfig zapcore.EncoderConfig
		var encoder zapcore.Encoder

		if format == "json" {
			encoderConfig = zap.NewProductionEncoderConfig()
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		} else {
			encoderConfig = zap.NewDevelopmentEncoderConfig()
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		}

		// Get output from config (stdout is default)
		output := cfg.GetString("logging.output")
		var writeSyncer zapcore.WriteSyncer

		if output == "stdout" || output == "" {
			writeSyncer = zapcore.AddSync(os.Stdout)
		} else {
			// For file output, open the file
			file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				// Fallback to stdout if file can't be opened
				writeSyncer = zapcore.AddSync(os.Stdout)
			} else {
				writeSyncer = zapcore.AddSync(file)
			}
		}

		// Create the core
		core := zapcore.NewCore(encoder, writeSyncer, level)

		// Create the logger
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	})
	return logger
}

// Sync flushes any buffered log entries
func Sync() error {
	if logger == nil {
		return nil
	}
	return logger.Sync()
}
