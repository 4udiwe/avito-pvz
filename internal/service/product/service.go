package product

import (
	"context"
	"errors"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/4udiwe/avito-pvz/internal/repository"
	"github.com/4udiwe/avito-pvz/pkg/transactor"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Service struct {
	productRepository   ProductsRepository
	receptionRepository ReceptionRepository
	txManager           transactor.Transactor
}

func New(p ProductsRepository, r ReceptionRepository, tx transactor.Transactor) *Service {
	return &Service{
		productRepository:   p,
		receptionRepository: r,
		txManager:           tx,
	}
}

func (s *Service) AddProduct(
	ctx context.Context,
	pointID uuid.UUID,
	productType entity.ProductType,
) (entity.Product, error) {
	logrus.Infof("Service: Adding product of type %s to point: %s", productType, pointID)
	var out entity.Product
	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Status check
		status, err := s.receptionRepository.GetLastReceptionStatus(ctx, pointID)

		if err != nil {
			logrus.Errorf("Service: Failed to get last reception status for point %s: %v", pointID, err)
			return err
		}

		if status != entity.ReceptionStatusInProgress {
			logrus.Warnf("Service: Reception already closed for point: %s", pointID)
			return ErrReceptionAlreadyClosed
		}

		// Create
		out, err = s.productRepository.Create(ctx, pointID, productType)
		return err
	})

	if err != nil {
		logrus.Errorf("Service: Failed to add product to point %s: %v", pointID, err)
		if errors.Is(err, repository.ErrNoPointFound) {
			return entity.Product{}, ErrNoPointFound
		}
		if errors.Is(err, repository.ErrNoReceptionFound) {
			return entity.Product{}, ErrNoReceptionFound
		}
		return entity.Product{}, err
	}

	logrus.Infof("Service: Product added: %+v", out)
	return out, nil
}

func (s *Service) DeleteLastProductFromReception(ctx context.Context, pointID uuid.UUID) error {
	logrus.Infof("Service: Deleting last product from reception for point: %s", pointID)
	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Status check
		status, err := s.receptionRepository.GetLastReceptionStatus(ctx, pointID)

		if err != nil {
			logrus.Errorf("Service: Failed to get last reception status for point %s: %v", pointID, err)
			return err
		}

		if status != entity.ReceptionStatusInProgress {
			logrus.Warnf("Service: Reception already closed for point: %s", pointID)
			return ErrReceptionAlreadyClosed
		}

		// Delete
		return s.productRepository.DeleteLastFromReception(ctx, pointID)
	})

	if err != nil {
		logrus.Errorf("Service: Failed to delete last product from reception for point %s: %v", pointID, err)
		if errors.Is(err, repository.ErrNoPointFound) {
			return ErrNoPointFound
		}
		if errors.Is(err, repository.ErrNoReceptionFound) {
			return ErrNoReceptionFound
		}
		return err
	}

	logrus.Infof("Service: Deleted last product from reception for point: %s", pointID)
	return nil
}
