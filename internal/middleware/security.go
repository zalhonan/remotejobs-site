package middleware

import (
	"net/http"
)

// Security - middleware для установки HTTP заголовков безопасности
type Security struct {
	// Флаг, указывающий, использует ли сайт HTTPS
	useHTTPS bool
}

// NewSecurity создает новый middleware для заголовков безопасности
func NewSecurity(useHTTPS bool) *Security {
	return &Security{
		useHTTPS: useHTTPS,
	}
}

// Middleware возвращает функцию middleware для установки заголовков безопасности
func (s *Security) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Content-Security-Policy - защита от XSS и инъекций
		// Разрешаем только ресурсы с нашего домена и CDN Bootstrap
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' https://cdn.jsdelivr.net; "+
				"style-src 'self' https://cdn.jsdelivr.net; "+
				"img-src 'self' data:; "+
				"font-src 'self' https://cdn.jsdelivr.net; "+
				"connect-src 'self'; "+
				"frame-ancestors 'none'")

		// X-Content-Type-Options - запрет на MIME-sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// X-Frame-Options - запрет на отображение в iframe (защита от кликджекинга)
		w.Header().Set("X-Frame-Options", "DENY")

		// Referrer-Policy - ограничение передачи Referer
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions-Policy - ограничение доступа к API браузера
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		// Strict-Transport-Security - принудительное использование HTTPS
		// Устанавливаем, только если сайт использует HTTPS
		if s.useHTTPS {
			// max-age=31536000 - один год в секундах
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		// Вызываем следующий обработчик
		next.ServeHTTP(w, r)
	})
}
