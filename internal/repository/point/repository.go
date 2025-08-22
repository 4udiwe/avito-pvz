package repo_point

import (
	"context"
	"errors"
	"fmt"

	"github.com/4udiwe/avito-pvz/internal/entity"
	repo "github.com/4udiwe/avito-pvz/internal/repository"
	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	*postgres.Postgres
}

func New(postgres *postgres.Postgres) *Repository {
	return &Repository{postgres}
}

func (r *Repository) Create(ctx context.Context, city string) (entity.Point, error) {
	logrus.Infof("Attempting to create point for city: %s", city)

	cityQuery := r.Builder.
		Select("id").
		From("cities").
		Where("name = ?", city)

	query, args, _ := r.Builder.
		Insert("points").
		Columns("city_id").
		Select(cityQuery).
		Suffix("RETURNING id, created_at").
		ToSql()

	point := entity.Point{City: city}
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(
		&point.ID,
		&point.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logrus.Warnf("No city found with name: %s", city)
			return entity.Point{}, repo.ErrNoCityFound
		}
		logrus.Errorf("Failed to create point for city %s: %v", city, err)
		return entity.Point{}, fmt.Errorf("PointRepository.Create - QueryRow: %w", err)
	}

	logrus.Infof("Point created: %+v", point)
	return point, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]entity.Point, error) {
	logrus.Info("Fetching all points")

	sql, args, _ := r.Builder.
		Select("points.id, created_at, cities.name AS city").
		From("points").
		InnerJoin("cities ON cities.id = points.city_id").
		ToSql()

	rows, err := r.GetTxManager(ctx).Query(ctx, sql, args...)
	if err != nil {
		logrus.Errorf("Failed to fetch points: %v", err)
		return nil, fmt.Errorf("PointRepository.GetAll - Query: %w", err)
	}
	defer rows.Close()

	var points []entity.Point
	for rows.Next() {
		var point entity.Point
		if err = rows.Scan(&point.ID, &point.CreatedAt, &point.City); err != nil {
			logrus.Errorf("Failed to scan point row: %v", err)
			return nil, fmt.Errorf("PointRepository.GetAll - rows.Scan: %w", err)
		}

		points = append(points, point)
	}

	if err = rows.Err(); err != nil {
		logrus.Errorf("Rows error after fetching points: %v", err)
		return nil, fmt.Errorf("PointRepository.GetAll - rows.Err: %w", err)
	}

	logrus.Infof("Fetched %d points", len(points))
	return points, nil
}
