package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zalhonan/remotejobs-site/internal/db"
	"github.com/zalhonan/remotejobs-site/internal/db/repository"
	"github.com/zalhonan/remotejobs-site/internal/domain/service"
	"github.com/zalhonan/remotejobs-site/internal/handler"
	"github.com/zalhonan/remotejobs-site/internal/logger"
	"github.com/zalhonan/remotejobs-site/internal/middleware"
	"github.com/zalhonan/remotejobs-site/internal/router"
	"go.uber.org/zap"
)

func main() {
	// Инициализация логгера
	appLogger, err := logger.InitLogger()
	if err != nil {
		panic("Cannot init logger: " + err.Error())
	}
	defer appLogger.Sync()

	appLogger.Info("Starting remote jobs website",
		zap.String("version", "1.0.0"),
	)

	ctx := context.Background()

	// Инициализация соединения с базой данных
	database, err := db.InitDB(ctx, appLogger)
	if err != nil {
		appLogger.Fatal("Не удалось инициализировать базу данных", zap.Error(err))
	}
	defer database.Close()

	// Создаем репозитории
	jobRepo := repository.NewJobRepository(database, appLogger)
	techRepo := repository.NewTechnologyRepository(database, appLogger)

	// Создаем сервисы
	jobService := service.NewJobService(jobRepo, techRepo, appLogger)
	technologyService := service.NewTechnologyService(techRepo, appLogger)

	// Создаем рендерер шаблонов
	templateRenderer, err := handler.NewTemplateRenderer("templates", "layout/base.html", appLogger)
	if err != nil {
		appLogger.Fatal("Не удалось инициализировать рендерер шаблонов", zap.Error(err))
	}

	// Создаем обработчики
	homeHandler := handler.NewHomeHandler(jobService, technologyService, templateRenderer, appLogger)
	jobHandler := handler.NewJobHandler(jobService, technologyService, templateRenderer, appLogger)

	// Создаем маршрутизатор
	appRouter := router.NewRouter(homeHandler, jobHandler, appLogger)

	// Создаем middleware для логирования запросов
	loggerMiddleware := middleware.NewLogger(appLogger)

	// Создаем middleware для заголовков безопасности
	// Устанавливаем useHTTPS в false, так как пока мы не используем HTTPS
	// Если сайт будет работать через HTTPS, нужно будет изменить на true
	securityMiddleware := middleware.NewSecurity(false)

	// Настройка HTTP-сервера
	server := &http.Server{
		Addr: ":8090", // Можно загружать из конфигурации
		// Применяем middleware по порядку: security -> logger -> router
		Handler:      securityMiddleware.Middleware(loggerMiddleware.Middleware(appRouter)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Запускаем сервер в отдельной горутине
	serverErrors := make(chan error, 1)
	go func() {
		appLogger.Info("Веб-сервер запущен", zap.String("address", server.Addr))
		serverErrors <- server.ListenAndServe()
	}()

	// Канал для сигналов завершения
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Ожидаем сигнал завершения или ошибку сервера
	select {
	case err := <-serverErrors:
		appLogger.Fatal("Ошибка сервера", zap.Error(err))

	case sig := <-shutdown:
		appLogger.Info("Получен сигнал завершения", zap.String("signal", sig.String()))

		// Даем время на корректное завершение
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Пытаемся корректно завершить сервер
		if err := server.Shutdown(ctx); err != nil {
			appLogger.Error("Ошибка при завершении сервера", zap.Error(err))
			if err := server.Close(); err != nil {
				appLogger.Fatal("Не удалось закрыть сервер", zap.Error(err))
			}
		}
	}

	appLogger.Info("Остановка вебсайта")
}
