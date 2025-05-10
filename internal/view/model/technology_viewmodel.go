package model

import (
	"github.com/zalhonan/remotejobs-site/internal/domain/entity"
)

// TechnologyViewModel модель представления для технологии
type TechnologyViewModel struct {
	ID        int64    // ID технологии
	Name      string   // Название технологии
	Keywords  []string // Ключевые слова
	SortOrder int      // Порядок сортировки
	URL       string   // URL для фильтрации вакансий по технологии
	JobsCount int64    // Количество доступных вакансий
}

// NewTechnologyViewModelFromEntity создает модель представления из доменной сущности
func NewTechnologyViewModelFromEntity(tech entity.Technology) TechnologyViewModel {
	// Формируем URL для фильтрации
	url := "/" + tech.Technology

	return TechnologyViewModel{
		ID:        tech.ID,
		Name:      tech.Technology,
		Keywords:  tech.Keywords,
		SortOrder: tech.SortOrder,
		URL:       url,
		JobsCount: tech.Count,
	}
}
