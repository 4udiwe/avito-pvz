package point

import (
	"context"
	"errors"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/4udiwe/avito-pvz/internal/repository"
	"github.com/4udiwe/avito-pvz/pkg/transactor"
	"github.com/sirupsen/logrus"
)

type Service struct {
	pointRepository PointRepository
	txManager       transactor.Transactor
}

func New(r PointRepository, tx transactor.Transactor) *Service {
	return &Service{
		pointRepository: r,
		txManager:       tx,
	}
}

func (s *Service) CreatePoint(ctx context.Context, city string) (entity.Point, error) {
	logrus.Infof("Service: Creating point for city: %s", city)
	point, err := s.pointRepository.Create(ctx, city)

	if err != nil {
		if errors.Is(err, repository.ErrNoCityFound) {
			logrus.Warnf("Service: No city found: %s", city)
			return entity.Point{}, ErrNoCityFound
		}
		logrus.Errorf("Service: Failed to create point for city %s: %v", city, err)
		return entity.Point{}, err
	}

	logrus.Infof("Service: Point created: %+v", point)
	return point, nil
}

func (s *Service) GetAllPoints(ctx context.Context) ([]entity.Point, error) {
	logrus.Info("Service: Fetching all points")
	points, err := s.pointRepository.GetAll(ctx)

	if err != nil {
		logrus.Errorf("Service: Failed to fetch all points: %v", err)
		return nil, err
	}

	logrus.Infof("Service: Fetched %d points", len(points))
	return points, nil
}
