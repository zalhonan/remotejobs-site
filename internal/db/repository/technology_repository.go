package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/zalhonan/remotejobs-site/internal/domain/entity"
	"go.uber.org/zap"
)

type TechnologyRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

// NewTechnologyRepository создает новый репозиторий для работы с технологиями
func NewTechnologyRepository(db *pgxpool.Pool, logger *zap.Logger) *TechnologyRepository {
	return &TechnologyRepository{
		db:     db,
		logger: logger,
	}
}

// GetAll возвращает все технологии, отсортированные по порядку сортировки
func (r *TechnologyRepository) GetAll(ctx context.Context) ([]entity.Technology, error) {
	query := `
		SELECT id, technology, keywords, sort_order
		FROM technologies
		ORDER BY sort_order DESC, technology ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список технологий: %w", err)
	}
	defer rows.Close()

	technologies := make([]entity.Technology, 0)
	for rows.Next() {
		var tech entity.Technology
		if err := rows.Scan(
			&tech.ID,
			&tech.Technology,
			&tech.Keywords,
			&tech.SortOrder,
		); err != nil {
			return nil, fmt.Errorf("не удалось обработать строку технологии: %w", err)
		}
		technologies = append(technologies, tech)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов запроса: %w", err)
	}

	return technologies, nil
}

// GetByName возвращает технологию по её имени
func (r *TechnologyRepository) GetByName(ctx context.Context, name string) (entity.Technology, error) {
	query := `
		SELECT id, technology, keywords, sort_order
		FROM technologies
		WHERE technology = $1
	`

	var tech entity.Technology
	err := r.db.QueryRow(ctx, query, name).Scan(
		&tech.ID,
		&tech.Technology,
		&tech.Keywords,
		&tech.SortOrder,
	)
	if err != nil {
		return entity.Technology{}, fmt.Errorf("не удалось получить технологию с именем=%s: %w", name, err)
	}

	return tech, nil
}

// Exists проверяет существование технологии по имени
func (r *TechnologyRepository) Exists(ctx context.Context, name string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM technologies WHERE technology = $1)"

	var exists bool
	err := r.db.QueryRow(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("не удалось проверить существование технологии %s: %w", name, err)
	}

	return exists, nil
}
