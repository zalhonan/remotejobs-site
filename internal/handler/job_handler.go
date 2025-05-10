package handler

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/zalhonan/remotejobs-site/internal/domain/service"
	"github.com/zalhonan/remotejobs-site/internal/view/model"
	"go.uber.org/zap"
)

type JobHandler struct {
	jobService        *service.JobService
	technologyService *service.TechnologyService
	templates         *TemplateRenderer
	logger            *zap.Logger
}

// NewJobHandler создает новый обработчик для страницы вакансии
func NewJobHandler(
	jobService *service.JobService,
	technologyService *service.TechnologyService,
	templates *TemplateRenderer,
	logger *zap.Logger,
) *JobHandler {
	return &JobHandler{
		jobService:        jobService,
		technologyService: technologyService,
		templates:         templates,
		logger:            logger,
	}
}

// Details обрабатывает запрос на страницу конкретной вакансии
func (h *JobHandler) Details(w http.ResponseWriter, r *http.Request, urlPath string) {
	// Регулярное выражение для извлечения слага из URL вида /job/slug
	re := regexp.MustCompile(`^/job/(.+)$`)
	matches := re.FindStringSubmatch(urlPath)

	if len(matches) != 2 {
		h.renderError(w, http.StatusNotFound, "Вакансия не найдена", "Запрошенная вакансия не существует или URL некорректен")
		return
	}

	// Извлекаем слаг из URL
	slug := matches[1]

	var jobID int64
	var err error

	// Проверяем, может ли слаг быть просто числом (ID)
	jobID, err = strconv.ParseInt(slug, 10, 64)
	if err != nil {
		// Если слаг не является просто числом, пытаемся извлечь ID из него
		// предполагая формат: ID-остальная-часть-слага
		idParts := strings.Split(slug, "-")
		if len(idParts) == 0 {
			h.renderError(w, http.StatusNotFound, "Вакансия не найдена", "Некорректный формат URL")
			return
		}

		jobID, err = strconv.ParseInt(idParts[0], 10, 64)
		if err != nil {
			h.logger.Error("Ошибка при преобразовании ID вакансии из слага",
				zap.Error(err),
				zap.String("slug", slug),
			)
			h.renderError(w, http.StatusBadRequest, "Неверный формат URL", "URL вакансии некорректен")
			return
		}
	}

	// Получаем вакансию по ID
	ctx := r.Context()
	job, err := h.jobService.GetByID(ctx, jobID)
	if err != nil {
		h.logger.Error("Ошибка при получении вакансии по ID",
			zap.Error(err),
			zap.Int64("jobId", jobID),
		)
		h.renderError(w, http.StatusNotFound, "Вакансия не найдена", "Запрошенная вакансия не существует или была удалена")
		return
	}

	// Преобразуем в view-модель
	jobViewModel := model.NewJobViewModelFromEntity(job, job.Slug)

	// Получаем похожие вакансии (например, с той же технологией)
	// Для простоты используем пустой массив
	relatedJobs := []model.JobViewModel{}

	// Получаем список всех технологий для меню
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

	// Формируем модель представления для детальной страницы вакансии
	viewModel := model.JobDetailViewModel{
		JobViewModel:    jobViewModel,
		RelatedJobs:     relatedJobs,
		PageTitle:       jobViewModel.Title,           // Используем заголовок вакансии в качестве заголовка страницы
		Technologies:    techViewModels,               // Добавляем список технологий для меню
		MetaDescription: jobViewModel.MetaDescription, // Используем мета-описание из модели вакансии
	}

	// Отображаем страницу
	if err := h.templates.Render(w, "pages/job_details.html", viewModel); err != nil {
		h.logger.Error("Ошибка при рендеринге шаблона job_details.html",
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
func (h *JobHandler) renderError(w http.ResponseWriter, statusCode int, title, message string) {
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
