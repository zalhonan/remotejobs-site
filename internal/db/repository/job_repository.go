package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/zalhonan/remotejobs-site/internal/domain/entity"
	"go.uber.org/zap"
)

type JobRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

// NewJobRepository создает новый репозиторий для работы с вакансиями
func NewJobRepository(db *pgxpool.Pool, logger *zap.Logger) *JobRepository {
	return &JobRepository{
		db:     db,
		logger: logger,
	}
}

// GetLatest возвращает последние вакансии с пагинацией
func (r *JobRepository) GetLatest(ctx context.Context, limit, offset int) ([]entity.JobRaw, error) {
	query := `
		SELECT id, content, title, source_link, main_technology, date_posted, date_parsed
		FROM jobs_raw
		ORDER BY date_posted DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список вакансий: %w", err)
	}
	defer rows.Close()

	jobs := make([]entity.JobRaw, 0)
	for rows.Next() {
		var job entity.JobRaw
		if err := rows.Scan(
			&job.ID,
			&job.Content,
			&job.Title,
			&job.SourceLink,
			&job.MainTechnology,
			&job.DatePosted,
			&job.DateParsed,
		); err != nil {
			return nil, fmt.Errorf("не удалось обработать строку вакансии: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов запроса: %w", err)
	}

	return jobs, nil
}

// GetByTechnology возвращает вакансии по конкретной технологии с пагинацией
func (r *JobRepository) GetByTechnology(ctx context.Context, technology string, limit, offset int) ([]entity.JobRaw, error) {
	query := `
		SELECT id, content, title, source_link, main_technology, date_posted, date_parsed
		FROM jobs_raw
		WHERE main_technology = $1
		ORDER BY date_posted DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, technology, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список вакансий по технологии %s: %w", technology, err)
	}
	defer rows.Close()

	jobs := make([]entity.JobRaw, 0)
	for rows.Next() {
		var job entity.JobRaw
		if err := rows.Scan(
			&job.ID,
			&job.Content,
			&job.Title,
			&job.SourceLink,
			&job.MainTechnology,
			&job.DatePosted,
			&job.DateParsed,
		); err != nil {
			return nil, fmt.Errorf("не удалось обработать строку вакансии: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов запроса: %w", err)
	}

	return jobs, nil
}

// GetByID возвращает вакансию по её ID
func (r *JobRepository) GetByID(ctx context.Context, id int64) (entity.JobRaw, error) {
	query := `
		SELECT id, content, title, source_link, main_technology, date_posted, date_parsed
		FROM jobs_raw
		WHERE id = $1
	`

	var job entity.JobRaw
	err := r.db.QueryRow(ctx, query, id).Scan(
		&job.ID,
		&job.Content,
		&job.Title,
		&job.SourceLink,
		&job.MainTechnology,
		&job.DatePosted,
		&job.DateParsed,
	)
	if err != nil {
		return entity.JobRaw{}, fmt.Errorf("не удалось получить вакансию с ID=%d: %w", id, err)
	}

	return job, nil
}

// GetTotalCount возвращает общее количество вакансий
func (r *JobRepository) GetTotalCount(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM jobs_raw"

	var count int
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("не удалось получить общее количество вакансий: %w", err)
	}

	return count, nil
}

// GetTotalCountByTechnology возвращает общее количество вакансий по технологии
func (r *JobRepository) GetTotalCountByTechnology(ctx context.Context, technology string) (int, error) {
	query := "SELECT COUNT(*) FROM jobs_raw WHERE main_technology = $1"

	var count int
	err := r.db.QueryRow(ctx, query, technology).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("не удалось получить количество вакансий по технологии %s: %w", technology, err)
	}

	return count, nil
}
