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
	pointRepository     PointRepository
	receptionRepository ReceptionRepository
	productRepository   ProductRepository
	txManager           transactor.Transactor
}

func New(pointRepo PointRepository,
	receptionRepo ReceptionRepository,
	productRepo ProductRepository,
	txManager transactor.Transactor,
) *Service {
	return &Service{
		pointRepository:     pointRepo,
		receptionRepository: receptionRepo,
		productRepository:   productRepo,
		txManager:           txManager,
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

func (s *Service) GetAllPointsFullInfo(ctx context.Context) ([]entity.PointFullInfo, error) {
	logrus.Info("Service: Fetching full info for all points")

	points, err := s.pointRepository.GetAll(ctx)
	if err != nil {
		logrus.Errorf("Service: Failed to fetch points: %v", err)
		return nil, err
	}

	var result []entity.PointFullInfo
	for _, point := range points {
		receptions, err := s.receptionRepository.GetAllByPoint(ctx, point.ID)
		if err != nil {
			logrus.Errorf("Service: Failed to fetch receptions for point %v: %v", point.ID, err)
			return nil, err
		}

		receptionsWithProducts := make([]entity.ReceptionWithProducts, 0)
		for _, reception := range receptions {
			products, err := s.productRepository.GetAllByReception(ctx, reception.ID)
			if err != nil {
				logrus.Errorf("Service: Failed to fetch products for reception %v: %v", reception.ID, err)
				return nil, err
			}
			receptionsWithProducts = append(receptionsWithProducts, entity.ReceptionWithProducts{
				Reception: reception,
				Products:  products,
			})
		}

		result = append(result, entity.PointFullInfo{
			Point:      point,
			Receptions: receptionsWithProducts,
		})
	}

	logrus.Infof("Service: Fetched full info for %d points", len(result))
	return result, nil
}
