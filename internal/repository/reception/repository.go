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
	"github.com/sirupsen/logrus"
)

type Repository struct {
	*postgres.Postgres
}

func New(postgres *postgres.Postgres) *Repository {
	return &Repository{postgres}
}

func (r *Repository) Open(ctx context.Context, pointID uuid.UUID) (entity.Reception, error) {
	logrus.Infof("Opening reception for point: %s", pointID)

	query, args, _ := r.Builder.
		Insert("receptions").
		Columns("point_id").
		Values(pointID).
		Suffix("RETURNING id, created_at, status").
		ToSql()

	reception := entity.Reception{PointID: pointID}
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(
		&reception.ID,
		&reception.CreatedAt,
		&reception.Status,
	)

	if err != nil {
		logrus.Errorf("Failed to open reception for point %s: %v", pointID, err)
		return entity.Reception{}, fmt.Errorf("ReceptionRepository.Open - create.Scan: %w", err)
	}

	logrus.Infof("Reception opened: %+v", reception)
	return reception, nil
}

func (r *Repository) GetLastReceptionStatus(ctx context.Context, pointID uuid.UUID) (entity.ReceptionStatus, error) {
	logrus.Infof("Fetching last reception status for point: %s", pointID)

	query, args, _ := r.Builder.
		Select("status").
		From("receptions").
		Where("point_id = ?", pointID).
		Limit(1).
		ToSql()

	var lastReceptionStatus entity.ReceptionStatus
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(&lastReceptionStatus)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logrus.Warnf("No receptions found for point: %s", pointID)
			return entity.ReceptionStatusClosed, nil
		}
		logrus.Errorf("Failed to fetch last reception status for point %s: %v", pointID, err)
		return "", fmt.Errorf("ReceptionRepository.GetLastReceptionStatus - lastStatus.Scan: %w", err)
	}

	logrus.Infof("Fetched last reception status for point %s: %s", pointID, lastReceptionStatus)
	return lastReceptionStatus, nil
}

func (r *Repository) GetLastReceptionProductsAmount(ctx context.Context, pointID uuid.UUID) (int, error) {
	logrus.Infof("Fetching last reception products amount for point: %s", pointID)

	query := `
        SELECT COUNT(p.id)
        FROM products p
        JOIN (
            SELECT id FROM receptions WHERE point_id = $1 ORDER BY created_at DESC LIMIT 1
        ) r ON p.reception_id = r.id
    `
	var productCount int
	err := r.GetTxManager(ctx).QueryRow(ctx, query, pointID).Scan(&productCount)
	if err != nil {
		logrus.Errorf("Failed to fetch product count for last reception of point %s: %v", pointID, err)
		return 0, fmt.Errorf("ReceptionRepository.GetLastReceptionProductsAmount - productCount.Scan: %w", err)
	}

	logrus.Infof("Fetched product count for last reception of point %s: %d", pointID, productCount)
	return productCount, nil
}

func (r *Repository) CloseLastReception(ctx context.Context, pointID uuid.UUID) error {
	logrus.Infof("Closing last reception for point: %s", pointID)

	query, args, _ := r.Builder.
		Update("receptions").
		Set("status", entity.ReceptionStatusClosed).
		Where("point_id = ? AND created_at = ("+
			"SELECT MAX(created_at) FROM receptions WHERE point_id = ?"+
			")", pointID, pointID).
		ToSql()

	result, err := r.GetTxManager(ctx).Exec(ctx, query, args...)
	if err != nil {
		logrus.Errorf("Failed to close last reception for point %s: %v", pointID, err)
		return fmt.Errorf("ReceptionRepository.CloseLastReception - Exec: %w", err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		logrus.Warnf("No reception found to close for point: %s", pointID)
		return repository.ErrNoReceptionFound
	}
	logrus.Infof("Closed last reception for point: %s", pointID)
	return nil
}

func (r *Repository) GetAllByPoint(ctx context.Context, pointID uuid.UUID) ([]entity.Reception, error) {
	logrus.Infof("Fetching all receptions for point: %s", pointID)

	query, args, _ := r.Builder.
		Select("id", "point_id", "created_at", "status").
		From("receptions").
		Where("point_id = ?", pointID).
		OrderBy("created_at ASC").
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, query, args...)
	if err != nil {
		logrus.Errorf("Failed to fetch receptions for point %s: %v", pointID, err)
		return nil, fmt.Errorf("ReceptionRepository.GetAllByPoint - Query: %w", err)
	}
	defer rows.Close()

	var receptions []entity.Reception
	for rows.Next() {
		var reception entity.Reception
		if err := rows.Scan(&reception.ID, &reception.PointID, &reception.CreatedAt, &reception.Status); err != nil {
			logrus.Errorf("Failed to scan reception row: %v", err)
			return nil, fmt.Errorf("ReceptionRepository.GetAllByPoint - Scan: %w", err)
		}
		receptions = append(receptions, reception)
	}
	if err := rows.Err(); err != nil {
		logrus.Errorf("Rows error after fetching receptions: %v", err)
		return nil, fmt.Errorf("ReceptionRepository.GetAllByPoint - rows.Err: %w", err)
	}

	logrus.Infof("Fetched %d receptions for point %s", len(receptions), pointID)
	return receptions, nil
}
