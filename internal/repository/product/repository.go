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
	"github.com/sirupsen/logrus"
)

type Repository struct {
	*postgres.Postgres
}

func New(postgres *postgres.Postgres) *Repository {
	return &Repository{postgres}
}

func (r *Repository) Create(ctx context.Context, pointID uuid.UUID, productType entity.ProductType) (entity.Product, error) {
	logrus.Infof("Attempting to create product of type %s for point: %s", productType, pointID)

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
			logrus.Warnf("No reception found for point: %s", pointID)
			return entity.Product{}, repository.ErrNoReceptionFound
		}
		logrus.Errorf("Failed to find reception for point %s: %v", pointID, err)
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
		logrus.Errorf("Failed to create product for reception %s: %v", receptionID, err)
		return entity.Product{}, fmt.Errorf("ProductRepository.Create - Scan: %w", err)
	}

	logrus.Infof("Product created: %+v", product)
	return product, nil
}

func (r *Repository) DeleteLastFromReception(ctx context.Context, pointID uuid.UUID) error {
	logrus.Infof("Deleting last product from reception for point: %s", pointID)

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
		logrus.Errorf("Failed to delete last product from reception for point %s: %v", pointID, err)
		return fmt.Errorf("ProductRepository.Delete - Exec: %w", err)
	}
	if result.RowsAffected() == 0 {
		logrus.Warnf("No product found to delete for point: %s", pointID)
		return repository.ErrNoProductFound
	}

	logrus.Infof("Deleted last product from reception for point: %s", pointID)
	return nil
}

func (r *Repository) GetAllByReception(ctx context.Context, receptionID uuid.UUID) ([]entity.Product, error) {
	logrus.Infof("Fetching all products for reception: %s", receptionID)

	query, args, _ := r.Builder.
		Select("id", "reception_id", "type", "created_at").
		From("products").
		Where("reception_id = ?", receptionID).
		OrderBy("created_at ASC").
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		logrus.Errorf("Failed to fetch products for reception %s: %v", receptionID, err)
		return nil, fmt.Errorf("ProductRepository.GetAllByReception - Query: %w", err)
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var product entity.Product
		if err := rows.Scan(&product.ID, &product.ReceptionID, &product.Type, &product.CreatedAt); err != nil {
			logrus.Errorf("Failed to scan product row: %v", err)
			return nil, fmt.Errorf("ProductRepository.GetAllByReception - Scan: %w", err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		logrus.Errorf("Rows error after fetching products: %v", err)
		return nil, fmt.Errorf("ProductRepository.GetAllByReception - rows.Err: %w", err)
	}

	logrus.Infof("Fetched %d products for reception %s", len(products), receptionID)
	return products, nil
}
