package service

import (
	"context"
	"regexp"
	"strings"

	"github.com/zalhonan/remotejobs-site/internal/db/repository"
	"github.com/zalhonan/remotejobs-site/internal/domain/entity"
	"go.uber.org/zap"
)

const (
	DefaultPageSize = 10
)

type JobService struct {
	jobRepo  *repository.JobRepository
	techRepo *repository.TechnologyRepository
	logger   *zap.Logger
}

// NewJobService создает новый сервис для работы с вакансиями
func NewJobService(
	jobRepo *repository.JobRepository,
	techRepo *repository.TechnologyRepository,
	logger *zap.Logger,
) *JobService {
	return &JobService{
		jobRepo:  jobRepo,
		techRepo: techRepo,
		logger:   logger,
	}
}

// GetLatest возвращает последние вакансии с пагинацией
func (s *JobService) GetLatest(ctx context.Context, page int) ([]entity.JobRaw, int, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * DefaultPageSize

	jobs, err := s.jobRepo.GetLatest(ctx, DefaultPageSize, offset)
	if err != nil {
		s.logger.Error("Не удалось получить последние вакансии",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("limit", DefaultPageSize),
			zap.Int("offset", offset),
		)
		return nil, 0, err
	}

	// Получаем общее количество страниц
	totalCount, err := s.jobRepo.GetTotalCount(ctx)
	if err != nil {
		s.logger.Error("Не удалось получить общее количество вакансий", zap.Error(err))
		return jobs, 0, nil
	}

	totalPages := (totalCount + DefaultPageSize - 1) / DefaultPageSize

	return jobs, totalPages, nil
}

// GetByTechnology возвращает вакансии по конкретной технологии с пагинацией
func (s *JobService) GetByTechnology(ctx context.Context, technology string, page int) ([]entity.JobRaw, int, error) {
	// Проверяем существование технологии
	exists, err := s.techRepo.Exists(ctx, technology)
	if err != nil {
		s.logger.Error("Ошибка при проверке существования технологии",
			zap.Error(err),
			zap.String("technology", technology),
		)
		return nil, 0, err
	}

	if !exists {
		s.logger.Warn("Запрошена несуществующая технология", zap.String("technology", technology))
		return []entity.JobRaw{}, 0, nil
	}

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * DefaultPageSize

	jobs, err := s.jobRepo.GetByTechnology(ctx, technology, DefaultPageSize, offset)
	if err != nil {
		s.logger.Error("Не удалось получить вакансии по технологии",
			zap.Error(err),
			zap.String("technology", technology),
			zap.Int("page", page),
			zap.Int("limit", DefaultPageSize),
			zap.Int("offset", offset),
		)
		return nil, 0, err
	}

	// Получаем общее количество страниц для этой технологии
	totalCount, err := s.jobRepo.GetTotalCountByTechnology(ctx, technology)
	if err != nil {
		s.logger.Error("Не удалось получить общее количество вакансий по технологии",
			zap.Error(err),
			zap.String("technology", technology),
		)
		return jobs, 0, nil
	}

	totalPages := (totalCount + DefaultPageSize - 1) / DefaultPageSize

	return jobs, totalPages, nil
}

// GetByID возвращает вакансию по её ID
func (s *JobService) GetByID(ctx context.Context, id int64) (entity.JobRaw, error) {
	job, err := s.jobRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Не удалось получить вакансию по ID",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return entity.JobRaw{}, err
	}

	return job, nil
}

// GenerateSlug генерирует slug из заголовка вакансии
func (s *JobService) GenerateSlug(title string) string {
	// Транслитерация кириллицы в латиницу могла бы быть здесь,
	// но для простоты просто заменяем все неалфавитно-цифровые символы на дефис
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	slug := reg.ReplaceAllString(strings.ToLower(title), "-")

	// Убираем начальные и конечные дефисы
	slug = strings.Trim(slug, "-")

	// Если slug пустой, возвращаем fallback
	if slug == "" {
		return "job"
	}

	return slug
}
