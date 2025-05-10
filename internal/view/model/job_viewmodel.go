package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/zalhonan/remotejobs-site/internal/domain/entity"
)

// JobViewModel модель представления для вакансии в списке
type JobViewModel struct {
	ID              int64     // ID вакансии
	Title           string    // Заголовок вакансии
	Content         string    // HTML содержимое вакансии для полной страницы
	ContentPreview  string    // Текстовое содержимое вакансии для превью
	SourceLink      string    // Ссылка на источник вакансии
	MainTechnology  string    // Основная технология
	DatePosted      time.Time // Дата публикации
	DatePostedStr   string    // Форматированная дата публикации
	Slug            string    // Часть URL для вакансии
	URL             string    // Полный URL вакансии
	MetaDescription string    // Мета-описание для SEO
}

// JobDetailViewModel модель представления для детальной страницы вакансии
type JobDetailViewModel struct {
	JobViewModel                          // Встраиваем базовую модель
	RelatedJobs     []JobViewModel        // Связанные вакансии
	PageTitle       string                // Заголовок страницы
	Technologies    []TechnologyViewModel // Список технологий для меню
	MetaDescription string                // Мета-описание для SEO
}

// JobListViewModel модель представления для списка вакансий
type JobListViewModel struct {
	Jobs            []JobViewModel        // Список вакансий
	Technologies    []TechnologyViewModel // Список технологий
	CurrentPage     int                   // Текущая страница
	TotalPages      int                   // Общее количество страниц
	Technology      string                // Текущая технология (если фильтр по технологии)
	IsFiltered      bool                  // Флаг, указывающий на наличие фильтра
	PrevPage        int                   // Предыдущая страница
	NextPage        int                   // Следующая страница
	PageTitle       string                // Заголовок страницы
	BaseURL         string                // Базовый URL для пагинации
	MetaDescription string                // Мета-описание для SEO
}

// createMetaDescriptionFromContent создает мета-описание из содержимого
func createMetaDescriptionFromContent(content string, technology string) string {
	// Очистка от HTML и экстра-пробелов
	trimmedContent := strings.TrimSpace(content)

	// Если контент пустой, возвращаем дефолтное описание с технологией
	if trimmedContent == "" {
		if technology != "" {
			return fmt.Sprintf("Удаленная вакансия по %s. Актуальные предложения о работе в IT с возможностью удаленной работы.", technology)
		}
		return "Актуальные удаленные вакансии в сфере IT. Работа из любой точки мира."
	}

	// Ограничиваем длину контента для мета-описания
	const maxDescriptionLength = 160
	if len(trimmedContent) > maxDescriptionLength {
		// Ищем последний пробел перед ограничением длины
		lastSpace := strings.LastIndex(trimmedContent[:maxDescriptionLength], " ")
		if lastSpace > 0 {
			trimmedContent = trimmedContent[:lastSpace]
		} else {
			trimmedContent = trimmedContent[:maxDescriptionLength]
		}
		trimmedContent += "..."
	}

	// Добавляем информацию о технологии, если она указана
	if technology != "" {
		return fmt.Sprintf("%s | Удаленная работа по %s", trimmedContent, technology)
	}

	return trimmedContent
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

	// Проверяем, не пустой ли slug
	if slug == "" {
		// Формируем простой fallback slug - просто числовой ID
		slug = fmt.Sprintf("%d", job.ID)
	}

	// Формируем URL вакансии, используя только слаг из базы данных
	url := fmt.Sprintf("/job/%s", slug)

	// Очищаем HTML-контент от лишних пробелов в начале
	content := strings.TrimLeft(job.Content, " \t\n\r")

	// Создаем мета-описание из ContentPure
	metaDescription := createMetaDescriptionFromContent(job.ContentPure, job.MainTechnology)

	return JobViewModel{
		ID:              job.ID,
		Title:           title,
		Content:         content,         // HTML содержимое для полной страницы
		ContentPreview:  job.ContentPure, // Текстовое содержимое для превью
		SourceLink:      job.SourceLink,
		MainTechnology:  job.MainTechnology,
		DatePosted:      job.DatePosted,
		DatePostedStr:   datePostedStr,
		Slug:            slug,
		URL:             url,
		MetaDescription: metaDescription,
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
	pageTitle := "Вакансии, удалённая работа в IT"

	// Создаем мета-описание для списка вакансий
	var metaDescription string
	if isFiltered {
		baseURL = "/" + technology + "/"
		pageTitle = "Вакансии по " + technology
		metaDescription = fmt.Sprintf("Актуальные удаленные вакансии по технологии %s. %d+ предложений о работе с возможностью работать из любой точки мира. Обновляется ежедневно.",
			technology, totalPages*10) // Примерная оценка количества вакансий
	} else {
		metaDescription = fmt.Sprintf("Свежие удаленные вакансии в IT. %d+ предложений о работе из любой точки мира. Фильтры по популярным технологиям, ежедневные обновления.",
			totalPages*10) // Примерная оценка количества вакансий
	}

	return JobListViewModel{
		Jobs:            jobs,
		Technologies:    technologies,
		CurrentPage:     currentPage,
		TotalPages:      totalPages,
		Technology:      technology,
		IsFiltered:      isFiltered,
		PrevPage:        prevPage,
		NextPage:        nextPage,
		PageTitle:       pageTitle,
		BaseURL:         baseURL,
		MetaDescription: metaDescription,
	}
}
