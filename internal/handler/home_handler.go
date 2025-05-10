package handler

import (
	"net/http"
	"strconv"

	"github.com/zalhonan/remotejobs-site/internal/domain/service"
	"github.com/zalhonan/remotejobs-site/internal/view/model"
	"go.uber.org/zap"
)

type HomeHandler struct {
	jobService        *service.JobService
	technologyService *service.TechnologyService
	templates         *TemplateRenderer
	logger            *zap.Logger
}

// NewHomeHandler создает новый обработчик для главной страницы
func NewHomeHandler(
	jobService *service.JobService,
	technologyService *service.TechnologyService,
	templates *TemplateRenderer,
	logger *zap.Logger,
) *HomeHandler {
	return &HomeHandler{
		jobService:        jobService,
		technologyService: technologyService,
		templates:         templates,
		logger:            logger,
	}
}

// Index обрабатывает запрос на главную страницу
func (h *HomeHandler) Index(w http.ResponseWriter, r *http.Request) {
	h.renderJobsList(w, r, "", 1)
}

// Page обрабатывает запрос на страницу с пагинацией
func (h *HomeHandler) Page(w http.ResponseWriter, r *http.Request, pageStr string) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		h.renderError(w, http.StatusBadRequest, "Неверный номер страницы", "Указанный номер страницы некорректен")
		return
	}

	h.renderJobsList(w, r, "", page)
}

// Technology обрабатывает запрос на страницу с вакансиями по технологии
func (h *HomeHandler) Technology(w http.ResponseWriter, r *http.Request, technology string) {
	h.renderJobsList(w, r, technology, 1)
}

// TechnologyPage обрабатывает запрос на страницу с пагинацией вакансий по технологии
func (h *HomeHandler) TechnologyPage(w http.ResponseWriter, r *http.Request, technology, pageStr string) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		h.renderError(w, http.StatusBadRequest, "Неверный номер страницы", "Указанный номер страницы некорректен")
		return
	}

	h.renderJobsList(w, r, technology, page)
}

// renderJobsList отображает список вакансий с учетом фильтров
func (h *HomeHandler) renderJobsList(w http.ResponseWriter, r *http.Request, technology string, page int) {
	ctx := r.Context()

	// Получаем список всех технологий
	technologies, err := h.technologyService.GetAll(ctx)
	if err != nil {
		h.logger.Error("Ошибка при получении списка технологий",
			zap.Error(err),
		)
		h.renderError(w, http.StatusInternalServerError, "Ошибка сервера", "Не удалось загрузить список технологий")
		return
	}

	// Преобразуем в view-модели для шаблона
	techViewModels := make([]model.TechnologyViewModel, 0, len(technologies))
	for _, tech := range technologies {
		techViewModel := model.NewTechnologyViewModelFromEntity(tech)
		techViewModels = append(techViewModels, techViewModel)
	}

	var jobs []model.JobViewModel
	var totalPages int

	// Получаем вакансии в зависимости от фильтра по технологии
	if technology != "" {
		// Проверяем, существует ли такая технология
		exists, err := h.technologyService.Exists(ctx, technology)
		if err != nil {
			h.logger.Error("Ошибка при проверке существования технологии",
				zap.Error(err),
				zap.String("technology", technology),
			)
			h.renderError(w, http.StatusInternalServerError, "Ошибка сервера", "Не удалось проверить существование технологии")
			return
		}

		if !exists {
			h.renderError(w, http.StatusNotFound, "Технология не найдена", "Запрошенная технология не существует")
			return
		}

		// Получаем вакансии по технологии
		jobsRaw, pages, err := h.jobService.GetByTechnology(ctx, technology, page)
		if err != nil {
			h.logger.Error("Ошибка при получении вакансий по технологии",
				zap.Error(err),
				zap.String("technology", technology),
				zap.Int("page", page),
			)
			h.renderError(w, http.StatusInternalServerError, "Ошибка сервера", "Не удалось загрузить вакансии по выбранной технологии")
			return
		}

		totalPages = pages

		// Преобразуем в view-модели
		jobs = make([]model.JobViewModel, 0, len(jobsRaw))
		for _, job := range jobsRaw {
			slug := h.jobService.GenerateSlug(job.Title)
			jobViewModel := model.NewJobViewModelFromEntity(job, slug)
			jobs = append(jobs, jobViewModel)
		}
	} else {
		// Получаем все вакансии
		jobsRaw, pages, err := h.jobService.GetLatest(ctx, page)
		if err != nil {
			h.logger.Error("Ошибка при получении последних вакансий",
				zap.Error(err),
				zap.Int("page", page),
			)
			h.renderError(w, http.StatusInternalServerError, "Ошибка сервера", "Не удалось загрузить список вакансий")
			return
		}

		totalPages = pages

		// Преобразуем в view-модели
		jobs = make([]model.JobViewModel, 0, len(jobsRaw))
		for _, job := range jobsRaw {
			slug := h.jobService.GenerateSlug(job.Title)
			jobViewModel := model.NewJobViewModelFromEntity(job, slug)
			jobs = append(jobs, jobViewModel)
		}
	}

	// Формируем модель представления для списка вакансий
	viewModel := model.NewJobListViewModel(jobs, techViewModels, page, totalPages, technology)

	// Отображаем страницу
	if err := h.templates.Render(w, "pages/home.html", viewModel); err != nil {
		h.logger.Error("Ошибка при рендеринге шаблона home.html",
			zap.Error(err),
		)
		// Если заголовок еще не был отправлен, устанавливаем код ошибки
		if w.Header().Get("Content-Type") == "" {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		} else {
			// Иначе просто пишем сообщение в тело ответа
			w.Write([]byte("Внутренняя ошибка сервера"))
		}
	}
}

// renderError отображает страницу с ошибкой
func (h *HomeHandler) renderError(w http.ResponseWriter, statusCode int, title, message string) {
	// Устанавливаем статус код только один раз
	w.WriteHeader(statusCode)

	viewModel := map[string]interface{}{
		"StatusCode":      statusCode,
		"Title":           title,
		"Message":         message,
		"PageTitle":       "Ошибка",
		"MetaDescription": "Ошибка на сайте удаленных вакансий в IT. " + message,
	}

	if err := h.templates.Render(w, "errors/error.html", viewModel); err != nil {
		h.logger.Error("Ошибка при рендеринге шаблона error.html",
			zap.Error(err),
		)
		// Не вызываем http.Error(), так как это вызовет повторную запись заголовка
		// Просто пишем сообщение об ошибке в тело ответа
		w.Write([]byte("Внутренняя ошибка сервера"))
	}
}
