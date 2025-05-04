package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Logger - middleware для логирования HTTP-запросов
type Logger struct {
	logger *zap.Logger
}

// NewLogger создает новый middleware для логирования
func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

// Middleware возвращает функцию middleware для логирования запросов
func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Оборачиваем ResponseWriter для получения статус-кода
		responseWriter := &ResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Вызываем следующий обработчик
		next.ServeHTTP(responseWriter, r)

		// Логируем запрос
		duration := time.Since(start)

		l.logger.Info("HTTP-запрос",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Int("status", responseWriter.statusCode),
			zap.Duration("duration", duration),
			zap.String("user_agent", r.UserAgent()),
		)
	})
}

// ResponseWriter - обертка для http.ResponseWriter для получения статус-кода
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader перехватывает статус-код
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
