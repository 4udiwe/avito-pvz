package reception

import (
	"context"
	"errors"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/4udiwe/avito-pvz/internal/repository"
	"github.com/4udiwe/avito-pvz/pkg/transactor"
	"github.com/google/uuid"
)

type Service struct {
	receptionRepository ReceptionRepository
	txManager       transactor.Transactor
}

func New(r ReceptionRepository, tx transactor.Transactor) *Service {
	return &Service{
		receptionRepository: r,
		txManager:       tx,
	}
}

func (s *Service) OpenReception(ctx context.Context, pointID uuid.UUID) (entity.Reception, error) {
	var reception entity.Reception
	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Status check
		status, err := s.receptionRepository.GetLastReceptionStatus(ctx, pointID)

		if err != nil {
			return err
		}

		if status != entity.ReceptionStatusClosed {
			return ErrLastReceptionNotClosed
		}

		// Open
		reception, err = s.receptionRepository.Open(ctx, pointID)

		return err
	})

	if err != nil {
		if errors.Is(err, repository.ErrNoPointFound) {
			return entity.Reception{}, ErrNoPointFound
		}
		return entity.Reception{}, err
	}

	return reception, nil
}

func (s *Service) CloseReception(ctx context.Context, pointID uuid.UUID) error {
	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Status check
		status, err := s.receptionRepository.GetLastReceptionStatus(ctx, pointID)

		if err != nil {
			return err
		}

		if status != entity.ReceptionStatusClosed {
			return ErrLastReceptionNotClosed
		}

		// Products amount check
		amount, err := s.receptionRepository.GetLastReceptionProductsAmount(ctx, pointID)

		if err != nil {
			return err
		}

		if amount == 0 {
			return ErrCannotCloseEmptyReception
		}

		// Close
		return s.receptionRepository.CloseLastReception(ctx, pointID)
	})

	if err != nil {
		if errors.Is(err, repository.ErrNoPointFound) {
			return ErrNoPointFound
		}
		if errors.Is(err, repository.ErrNoReceptionFound) {
			return ErrNoReceptionFound
		}
		return err
	}

	return nil
}
