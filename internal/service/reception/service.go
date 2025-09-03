package reception

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
	receptionRepository ReceptionRepository
	txManager           transactor.Transactor
	metrics             Metrics
}

func New(r ReceptionRepository, tx transactor.Transactor, m Metrics) *Service {
	return &Service{
		receptionRepository: r,
		txManager:           tx,
		metrics:             m,
	}
}

func (s *Service) OpenReception(ctx context.Context, pointID uuid.UUID) (entity.Reception, error) {
	logrus.Infof("Service: Opening reception for point: %s", pointID)
	var reception entity.Reception
	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Point existence check
		exists, err := s.receptionRepository.CheckIfPointExists(ctx, pointID)
		if err != nil {
			logrus.Errorf("Service: Failed to check if point exists for point %s: %v", pointID, err)
			return err
		}
		if !exists {
			logrus.Warnf("Service: Point does not exist: %s", pointID)
			return ErrNoPointFound
		}

		// Status check
		status, err := s.receptionRepository.GetLastReceptionStatus(ctx, pointID)

		if err != nil {
			logrus.Errorf("Service: Failed to get last reception status for point %s: %v", pointID, err)
			return err
		}

		if status != entity.ReceptionStatusClosed {
			logrus.Warnf("Service: Last reception not closed for point: %s", pointID)
			return ErrLastReceptionNotClosed
		}

		// Open
		reception, err = s.receptionRepository.Open(ctx, pointID)

		return err
	})

	if err != nil {
		logrus.Errorf("Service: Failed to open reception for point %s: %v", pointID, err)
		if !errors.Is(err, ErrNoPointFound) && !errors.Is(err, ErrLastReceptionNotClosed) {
			s.metrics.ErrInc()
		}
		return entity.Reception{}, err
	} 

	logrus.Infof("Service: Reception opened: %+v", reception)
	s.metrics.Inc()
	return reception, nil
}

func (s *Service) CloseReception(ctx context.Context, pointID uuid.UUID) error {
	logrus.Infof("Service: Closing reception for point: %s", pointID)
	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Point existence check
		exists, err := s.receptionRepository.CheckIfPointExists(ctx, pointID)
		if err != nil {
			logrus.Errorf("Service: Failed to check if point exists for point %s: %v", pointID, err)
			return err
		}
		if !exists {
			logrus.Warnf("Service: Point does not exist: %s", pointID)
			return ErrNoPointFound
		}

		// Status check
		status, err := s.receptionRepository.GetLastReceptionStatus(ctx, pointID)

		if err != nil {
			logrus.Errorf("Service: Failed to get last reception status for point %s: %v", pointID, err)
			return err
		}

		if status == entity.ReceptionStatusClosed {
			logrus.Warnf("Service: Last reception already closed for point: %s", pointID)
			return ErrLastReceptionAlreadyClosed
		}

		// Products amount check
		amount, err := s.receptionRepository.GetLastReceptionProductsAmount(ctx, pointID)

		if err != nil {
			logrus.Errorf("Service: Failed to get products amount for point %s: %v", pointID, err)
			return err
		}

		if amount == 0 {
			logrus.Warnf("Service: Cannot close empty reception for point: %s", pointID)
			return ErrCannotCloseEmptyReception
		}

		// Close
		return s.receptionRepository.CloseLastReception(ctx, pointID)
	})

	if err != nil {
		if errors.Is(err, repository.ErrNoReceptionFound) {
			return ErrNoReceptionFound
		}
		logrus.Errorf("Service: Failed to close reception for point %s: %v", pointID, err)
		return err
	}

	logrus.Infof("Service: Reception closed for point: %s", pointID)
	return nil
}
