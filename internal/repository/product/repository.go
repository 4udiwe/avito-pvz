package repo_product

import (
	"context"
	"errors"
	"fmt"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/4udiwe/avito-pvz/internal/repository"
	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	*postgres.Postgres
}

func New(postgres *postgres.Postgres) *Repository {
	return &Repository{postgres}
}

func (r *Repository) Create(ctx context.Context, pointID uuid.UUID, productType entity.ProductType) (entity.Product, error) {
	query, args, _ := r.Builder.
		Select("id").
		From("receptions").
		Where(squirrel.Eq{"point_id": pointID}).
		OrderBy("created_at DESC").
		Limit(1).
		ToSql()

	var receptionID uuid.UUID
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(&receptionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Product{}, repository.ErrNoReceptionFound
		}
		return entity.Product{}, fmt.Errorf("ProductRepository.Create - find reception: %w", err)
	}

	query, args, _ = r.Builder.
		Insert("products").
		Columns("reception_id", "type").
		Values(receptionID, productType).
		Suffix("RETURNING id, created_at").
		ToSql()

	product := entity.Product{
		ReceptionID: receptionID,
		Type:        productType,
	}
	err = r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(&product.ID, &product.CreatedAt)

	if err != nil {
		return entity.Product{}, fmt.Errorf("ProductRepository.Create - Scan: %w", err)
	}

	return product, nil
}

func (r *Repository) DeleteLastFromReception(ctx context.Context, pointID uuid.UUID) error {
	query, args, _ := r.Builder.
		Delete("products").
		Where("id = ("+
			"SELECT p.id FROM products p "+
			"JOIN receptions r ON p.reception_id = r.id "+
			"WHERE r.point_id = ? "+
			"ORDER BY p.created_at DESC, p.id DESC "+
			"LIMIT 1"+
			")", pointID).
		ToSql()

	result, err := r.GetTxManager(ctx).Exec(ctx, query, args...)

	if err != nil {
		return fmt.Errorf("ProductRepository.Delete - Exec: %w", err)
	}
	if result.RowsAffected() == 0 {
		return repository.ErrNoProductFound
	}

	return nil
}
