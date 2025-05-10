package router

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/zalhonan/remotejobs-site/internal/handler"
	"go.uber.org/zap"
)

// Router обрабатывает HTTP-запросы и направляет их соответствующим обработчикам
type Router struct {
	homeHandler *handler.HomeHandler
	jobHandler  *handler.JobHandler
	logger      *zap.Logger
}

// NewRouter создает новый маршрутизатор
func NewRouter(
	homeHandler *handler.HomeHandler,
	jobHandler *handler.JobHandler,
	logger *zap.Logger,
) *Router {
	return &Router{
		homeHandler: homeHandler,
		jobHandler:  jobHandler,
		logger:      logger,
	}
}

// ServeHTTP обрабатывает HTTP-запросы
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	// Статические файлы
	if match, _ := regexp.MatchString(`^/static/`, path); match {
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, req)
		return
	}

	// Главная страница
	if path == "/" {
		r.homeHandler.Index(w, req)
		return
	}

	// Страница вакансии
	if match, _ := regexp.MatchString(`^/job/\d+`, path); match {
		r.jobHandler.Details(w, req, path)
		return
	}

	// Пагинация на главной странице
	if match, _ := regexp.MatchString(`^/\d+$`, path); match {
		pageStr := path[1:] // Удаляем начальный слеш
		r.homeHandler.Page(w, req, pageStr)
		return
	}

	// Проверяем, что технология не число (чтобы не конфликтовало с пагинацией)
	reIsNumber := regexp.MustCompile(`^/\d+$`)
	if !reIsNumber.MatchString(path) {
		// Пагинация для технологии
		reTechPage := regexp.MustCompile(`^/([^/]+)/(\d+)$`)
		if matches := reTechPage.FindStringSubmatch(path); len(matches) == 3 {
			technology := matches[1]
			pageStr := matches[2]
			r.homeHandler.TechnologyPage(w, req, technology, pageStr)
			return
		}

		// Страница технологии
		reTech := regexp.MustCompile(`^/([^/]+)$`)
		if matches := reTech.FindStringSubmatch(path); len(matches) == 2 {
			technology := matches[1]
			// Проверяем, что это не номер страницы
			if _, err := strconv.Atoi(technology); err != nil {
				r.homeHandler.Technology(w, req, technology)
				return
			}
		}
	}

	// Если ни один маршрут не подошел, возвращаем 404
	http.NotFound(w, req)
}

// Routes возвращает обработчик HTTP-запросов
func (r *Router) Routes() http.Handler {
	return r
}
