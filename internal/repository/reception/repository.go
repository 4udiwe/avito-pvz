package repo_reception

import (
	"context"
	"errors"
	"fmt"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/4udiwe/avito-pvz/internal/repository"
	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	*postgres.Postgres
}

func New(postgres *postgres.Postgres) *Repository {
	return &Repository{postgres}
}

func (r *Repository) Open(ctx context.Context, pointID uuid.UUID) (entity.Reception, error) {
	query, args, _ := r.Builder.
		Insert("receptions").
		Columns("point_id").
		Suffix("RETURNING id, created_at, status").
		ToSql()

	reception := entity.Reception{PointID: pointID}
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(
		&reception.ID,
		&reception.CreatedAt,
		&reception.Status,
	)

	if err != nil {
		return entity.Reception{}, fmt.Errorf("ReceptionRepository.Create - create.Scan: %w", err)
	}

	return reception, nil
}

func (r *Repository) GetLastReceptionStatus(ctx context.Context, pointID uuid.UUID) (entity.ReceptionStatus, error) {
	query, args, _ := r.Builder.
		Select("status").
		From("receptions").
		Where("point_id = ?", pointID).
		OrderBy("created_at DESC").
		Limit(1).
		ToSql()

	var lastReceptionStatus entity.ReceptionStatus
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(&lastReceptionStatus)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repository.ErrNoPointFound
		}
		return "", fmt.Errorf("ReceptionRepository.Create - lastStatus.Scan: %w", err)
	}

	return lastReceptionStatus, nil
}

func (r *Repository) GetLastReceptionProductsAmount(ctx context.Context, pointID uuid.UUID) (int, error) {
	subquery := r.Builder.
		Select("id").
		From("receptions").
		Where("point_id = ?", pointID).
		OrderBy("created_at DESC").
		Limit(1)

	query, args, _ := r.Builder.
		Select("COUNT(p.id)").
		From("products p").
		JoinClause("JOIN (?) AS r ON p.reception_id = r.id", subquery).
		ToSql()

	var productCount int
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(&productCount)
	if err != nil {
		return 0, fmt.Errorf("ReceptionRepository.GetLastReceptionProductsAmount - productCount.Scan: %w", err)
	}

	return productCount, nil
}

func (r *Repository) CloseLastReception(ctx context.Context, pointID uuid.UUID) error {
	query, args, _ := r.Builder.
		Update("receptions").
		Set("status", entity.ReceptionStatusClosed).
		Where("point_id = ? AND created_at = ("+
			"SELECT MAX(created_at) FROM receptions WHERE point_id = ?"+
			")", pointID, pointID).
		ToSql()

	result, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("ReceptionRepository.CloseLastReception - Exec: %w", err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return repository.ErrNoReceptionFound
	}
	return nil
}
