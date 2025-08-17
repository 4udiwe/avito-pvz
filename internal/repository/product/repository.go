package repo_product

import (
	"context"
	"fmt"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/4udiwe/avito-pvz/internal/repository"
	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/google/uuid"
)

type Repository struct {
	*postgres.Postgres
}

func New(postgres *postgres.Postgres) *Repository {
	return &Repository{postgres}
}

func (r *Repository) Create(ctx context.Context, p createProduct) (entity.Product, error) {
	query, args, _ := r.Builder.
		Insert("products").
		Columns("reception_id", "type").
		Values(p.ReceptionID, p.Type).
		Suffix("RETURNING id, created_at").
		ToSql()

	product := entity.Product{
		ReceptionID: p.ReceptionID,
		Type:        p.Type,
	}
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(&product.ID, &product.CreatedAt)

	if err != nil {
		return entity.Product{}, fmt.Errorf("ProductRepository.Create - Scan: %w", err)
	}

	return product, nil
}

func (r *Repository) DeleteById(ctx context.Context, productID uuid.UUID) error {
	query, args, _ := r.Builder.
		Delete("products").
		Where("id = ?", productID).
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
