package db

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Приоритет имеют переменные окружения
func InitDB(ctx context.Context, logger *zap.Logger) (*pgxpool.Pool, error) {
	// Загружаем .env файл, но не останавливаем выполнение, если его нет
	if err := godotenv.Load(".env"); err != nil {
		logger.Warn("Не удалось загрузить .env файл из корня проекта", zap.Error(err))
	} else {
		logger.Info("Конфигурация .env успешно загружена")
	}

	// Формируем строку подключения к базе данных
	dbHost := os.Getenv("PG_HOST")
	dbPort := os.Getenv("PG_PORT")
	dbName := os.Getenv("PG_DATABASE_NAME")
	dbUser := os.Getenv("PG_USER")
	dbPassword := os.Getenv("PG_PASSWORD")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	dbParams := []string{}
	if dbHost != "" {
		dbParams = append(dbParams, fmt.Sprintf("host=%s", dbHost))
	}
	if dbPort != "" {
		dbParams = append(dbParams, fmt.Sprintf("port=%s", dbPort))
	}
	if dbName != "" {
		dbParams = append(dbParams, fmt.Sprintf("dbname=%s", dbName))
	}
	if dbUser != "" {
		dbParams = append(dbParams, fmt.Sprintf("user=%s", dbUser))
	}
	if dbPassword != "" {
		dbParams = append(dbParams, fmt.Sprintf("password=%s", dbPassword))
	}
	if dbSSLMode != "" {
		dbParams = append(dbParams, fmt.Sprintf("sslmode=%s", dbSSLMode))
	}

	dbDSN := strings.Join(dbParams, " ")

	if dbDSN == "" {
		return nil, fmt.Errorf("не удалось создать строку подключения к БД: не найдены необходимые переменные окружения")
	}

	// Подключаемся к БД
	poolConfig, err := pgxpool.ParseConfig(dbDSN)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить строку подключения к базе данных: %w", err)
	}

	db, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	logger.Info("Успешное подключение к базе данных")

	return db, nil
}
