package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/zalhonan/remotejobs-site/internal/domain/entity"
)

// JobViewModel модель представления для вакансии в списке
type JobViewModel struct {
	ID             int64     // ID вакансии
	Title          string    // Заголовок вакансии
	Content        string    // HTML содержимое вакансии для полной страницы
	ContentPreview string    // Текстовое содержимое вакансии для превью
	SourceLink     string    // Ссылка на источник вакансии
	MainTechnology string    // Основная технология
	DatePosted     time.Time // Дата публикации
	DatePostedStr  string    // Форматированная дата публикации
	Slug           string    // Часть URL для вакансии
	URL            string    // Полный URL вакансии
}

// JobDetailViewModel модель представления для детальной страницы вакансии
type JobDetailViewModel struct {
	JobViewModel                       // Встраиваем базовую модель
	RelatedJobs  []JobViewModel        // Связанные вакансии
	PageTitle    string                // Заголовок страницы
	Technologies []TechnologyViewModel // Список технологий для меню
}

// JobListViewModel модель представления для списка вакансий
type JobListViewModel struct {
	Jobs         []JobViewModel        // Список вакансий
	Technologies []TechnologyViewModel // Список технологий
	CurrentPage  int                   // Текущая страница
	TotalPages   int                   // Общее количество страниц
	Technology   string                // Текущая технология (если фильтр по технологии)
	IsFiltered   bool                  // Флаг, указывающий на наличие фильтра
	PrevPage     int                   // Предыдущая страница
	NextPage     int                   // Следующая страница
	PageTitle    string                // Заголовок страницы
	BaseURL      string                // Базовый URL для пагинации
}

// NewJobViewModelFromEntity создает модель представления из доменной сущности
func NewJobViewModelFromEntity(job entity.JobRaw, slug string) JobViewModel {
	// Форматируем дату для отображения
	datePostedStr := job.DatePosted.Format("02.01.2006")

	// Используем поле Title из JobRaw
	title := job.Title
	// Если title не задан, формируем его из main_technology
	if title == "" {
		title = "Вакансия по " + job.MainTechnology
	}

	// Формируем URL вакансии
	url := fmt.Sprintf("/job/%d-%s", job.ID, slug)

	// Очищаем HTML-контент от лишних пробелов в начале
	content := strings.TrimLeft(job.Content, " \t\n\r")

	return JobViewModel{
		ID:             job.ID,
		Title:          title,
		Content:        content,         // HTML содержимое для полной страницы
		ContentPreview: job.ContentPure, // Текстовое содержимое для превью
		SourceLink:     job.SourceLink,
		MainTechnology: job.MainTechnology,
		DatePosted:     job.DatePosted,
		DatePostedStr:  datePostedStr,
		Slug:           slug,
		URL:            url,
	}
}

// NewJobListViewModel создает модель представления списка вакансий
func NewJobListViewModel(
	jobs []JobViewModel,
	technologies []TechnologyViewModel,
	currentPage, totalPages int,
	technology string,
) JobListViewModel {
	isFiltered := technology != ""
	prevPage := currentPage - 1
	if prevPage < 1 {
		prevPage = 1
	}

	nextPage := currentPage + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	baseURL := "/"
	pageTitle := "Все вакансии"

	if isFiltered {
		baseURL = "/" + technology + "/"
		pageTitle = "Вакансии по " + technology
	}

	return JobListViewModel{
		Jobs:         jobs,
		Technologies: technologies,
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
		Technology:   technology,
		IsFiltered:   isFiltered,
		PrevPage:     prevPage,
		NextPage:     nextPage,
		PageTitle:    pageTitle,
		BaseURL:      baseURL,
	}
}
