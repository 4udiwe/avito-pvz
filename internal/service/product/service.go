package product

import (
	"context"
	"errors"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/4udiwe/avito-pvz/internal/repository"
	repo_product "github.com/4udiwe/avito-pvz/internal/repository/product"
	"github.com/4udiwe/avito-pvz/pkg/transactor"
	"github.com/google/uuid"
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
	product repo_product.CreateProduct,
) (entity.Product, error) {
	var out entity.Product
	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Status check
		status, err := s.receptionRepository.GetLastReceptionStatus(ctx, pointID)

		if err != nil {
			return err
		}

		if status != entity.ReceptionStatusInProgress {
			return ErrReceptionAlreadyClosed
		}

		// Create
		out, err = s.productRepository.Create(ctx, product)
		return err
	})

	if err != nil {
		if errors.Is(err, repository.ErrNoPointFound) {
			return entity.Product{}, ErrNoPointFound
		}
		return entity.Product{}, err
	}

	return out, nil
}

func (s *Service) DeleteLastProductFromReception(ctx context.Context, pointID uuid.UUID) error {
	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Status check
		status, err := s.receptionRepository.GetLastReceptionStatus(ctx, pointID)

		if err != nil {
			return err
		}

		if status != entity.ReceptionStatusInProgress {
			return ErrReceptionAlreadyClosed
		}

		// Delete
		return s.productRepository.DeleteLastFromReception(ctx, pointID)
	})

	if err != nil {
		if errors.Is(err, repository.ErrNoPointFound) {
			return ErrNoPointFound
		}
		return err
	}

	return nil
}
