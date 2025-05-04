package service

import (
	"context"

	"github.com/zalhonan/remotejobs-site/internal/db/repository"
	"github.com/zalhonan/remotejobs-site/internal/domain/entity"
	"go.uber.org/zap"
)

type TechnologyService struct {
	techRepo *repository.TechnologyRepository
	logger   *zap.Logger
}

// NewTechnologyService создает новый сервис для работы с технологиями
func NewTechnologyService(techRepo *repository.TechnologyRepository, logger *zap.Logger) *TechnologyService {
	return &TechnologyService{
		techRepo: techRepo,
		logger:   logger,
	}
}

// GetAll возвращает все технологии, отсортированные по порядку сортировки
func (s *TechnologyService) GetAll(ctx context.Context) ([]entity.Technology, error) {
	technologies, err := s.techRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error("Не удалось получить список технологий", zap.Error(err))
		return nil, err
	}

	return technologies, nil
}

// GetByName возвращает технологию по её имени
func (s *TechnologyService) GetByName(ctx context.Context, name string) (entity.Technology, error) {
	technology, err := s.techRepo.GetByName(ctx, name)
	if err != nil {
		s.logger.Error("Не удалось получить технологию по имени",
			zap.Error(err),
			zap.String("name", name),
		)
		return entity.Technology{}, err
	}

	return technology, nil
}

// Exists проверяет существование технологии по имени
func (s *TechnologyService) Exists(ctx context.Context, name string) (bool, error) {
	exists, err := s.techRepo.Exists(ctx, name)
	if err != nil {
		s.logger.Error("Ошибка при проверке существования технологии",
			zap.Error(err),
			zap.String("name", name),
		)
		return false, err
	}

	return exists, nil
}
