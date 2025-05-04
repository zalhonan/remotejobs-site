package logger

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// InitLogger создает и возвращает новый экземпляр логгера
func InitLogger() (*zap.Logger, error) {
	// Загружаем переменные окружения из .env, если файл существует
	// Игнорируем ошибку, если файл не найден
	workDir, _ := os.Getwd()
	_ = godotenv.Load(filepath.Join(workDir, ".env"))

	// Получаем параметры из переменных окружения или используем значения по умолчанию
	betterStackKey := os.Getenv("BETTERSTACK_KEY")

	betterStackURL := os.Getenv("BETTERSTACK_URL")

	logger, err := NewLogger(
		betterStackKey,
		betterStackURL,
		50,
		10,
	)

	return logger, err
}
