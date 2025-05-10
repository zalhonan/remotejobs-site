package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zalhonan/remotejobs-site/internal/handler"
	"go.uber.org/zap"
)

// NewRouter создает новый маршрутизатор на основе Chi
func NewRouter(
	homeHandler *handler.HomeHandler,
	jobHandler *handler.JobHandler,
	logger *zap.Logger,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Timeout(60 * time.Second))

	// Статические файлы
	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Маршруты
	r.Get("/", homeHandler.Index)

	// Пагинация на главной странице
	r.Get("/{page}", func(w http.ResponseWriter, r *http.Request) {
		page := chi.URLParam(r, "page")
		// Проверяем, что это номер, а не технология
		if _, err := strconv.Atoi(page); err == nil {
			homeHandler.Page(w, r, page)
			return
		}
		// Если это не номер, значит это технология
		homeHandler.Technology(w, r, page)
	})

	// Пагинация для технологии
	r.Get("/{technology}/{page}", func(w http.ResponseWriter, r *http.Request) {
		technology := chi.URLParam(r, "technology")
		page := chi.URLParam(r, "page")
		homeHandler.TechnologyPage(w, r, technology, page)
	})

	// Страница вакансии
	r.Get("/job/{jobID}-{slug}", func(w http.ResponseWriter, r *http.Request) {
		// Chi уже разбирает URL за нас
		path := "/job/" + chi.URLParam(r, "jobID") + "-" + chi.URLParam(r, "slug")
		jobHandler.Details(w, r, path)
	})

	return r
}
