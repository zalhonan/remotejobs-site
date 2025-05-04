package logger

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(
	betterStackToken string,
	betterStackURL string,
	batchSize int,
	flushInterval time.Duration,

) (*zap.Logger, error) {
	// Using relative path to reach project root from internal/logger
	logDir := filepath.Join("..", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	betterStackSink := NewBetterStackSink(
		betterStackToken,
		betterStackURL,
		batchSize,     // batch size - send logs in batches of 10
		flushInterval, // flush interval - also flush every 30 seconds
	)

	// Create BetterStack-specific encoder config
	betterStackEncoderConfig := zapcore.EncoderConfig{
		TimeKey:     "dt",
		MessageKey:  "message",
		LevelKey:    "level",
		LineEnding:  zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime: zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format("2006-01-02 15:04:05 UTC"))
		}),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Configure standard encoders
	productionEncoderConfig := zap.NewProductionEncoderConfig()
	productionEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	debugEncoderConfig := zap.NewProductionEncoderConfig()
	debugEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	debugEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Create encoders
	jsonEncoder := zapcore.NewJSONEncoder(productionEncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(debugEncoderConfig)

	// Prepare log file
	logFile := filepath.Join(logDir, "app.log")
	fileWriter, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Create multi-output core
	core := zapcore.NewTee(
		// File output (JSON format)
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(fileWriter), zap.InfoLevel),
		// Console output (with colors)
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
		// BetterStack output (custom JSON format)
		zapcore.NewCore(zapcore.NewJSONEncoder(betterStackEncoderConfig), zapcore.AddSync(betterStackSink), zap.InfoLevel),
	)

	logger := zap.New(core, zap.AddStacktrace(zap.ErrorLevel))
	return logger, nil
}
