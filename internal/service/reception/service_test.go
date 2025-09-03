package reception_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4udiwe/avito-pvz/internal/entity"
	mock_transactor "github.com/4udiwe/avito-pvz/internal/mocks"
	"github.com/4udiwe/avito-pvz/internal/repository"
	service "github.com/4udiwe/avito-pvz/internal/service/reception"
	mock_reception "github.com/4udiwe/avito-pvz/internal/service/reception/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOpenReception(t *testing.T) {
	var (
		ctx          = context.Background()
		pointID      = uuid.New()
		arbitraryErr = errors.New("arbitraryErr")

		emptyStatus entity.ReceptionStatus = ""
	)

	reception := entity.Reception{
		ID:      uuid.New(),
		PointID: pointID,
		Status:  entity.ReceptionStatusInProgress,
	}

	type MockBehavior func(
		r *mock_reception.MockReceptionRepository,
		t *mock_transactor.MockTransactor,
		m *mock_reception.MockMetrics,
	)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		want         entity.Reception
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(r *mock_reception.MockReceptionRepository, t *mock_transactor.MockTransactor, m *mock_reception.MockMetrics) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				r.EXPECT().CheckIfPointExists(ctx, pointID).Return(true, nil).Times(1)
				r.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusClosed, nil).Times(1)
				r.EXPECT().Open(ctx, pointID).Return(reception, nil).Times(1)
				m.EXPECT().Inc().Times(1)
			},
			want:    reception,
			wantErr: nil,
		},
		{
			name: "failed to get status",
			mockBehavior: func(r *mock_reception.MockReceptionRepository, t *mock_transactor.MockTransactor, m *mock_reception.MockMetrics) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				r.EXPECT().CheckIfPointExists(ctx, pointID).Return(true, nil).Times(1)
				r.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(emptyStatus, arbitraryErr).Times(1)
				m.EXPECT().ErrInc().Times(1)
			},
			want:    entity.Reception{},
			wantErr: arbitraryErr,
		},
		{
			name: "last reception not closed",
			mockBehavior: func(r *mock_reception.MockReceptionRepository, t *mock_transactor.MockTransactor, m *mock_reception.MockMetrics) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				r.EXPECT().CheckIfPointExists(ctx, pointID).Return(true, nil).Times(1)
				r.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusInProgress, nil).Times(1)
			},
			want:    entity.Reception{},
			wantErr: service.ErrLastReceptionNotClosed,
		},
		{
			name: "no point found",
			mockBehavior: func(r *mock_reception.MockReceptionRepository, t *mock_transactor.MockTransactor, m *mock_reception.MockMetrics) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				r.EXPECT().CheckIfPointExists(ctx, pointID).Return(false, nil).Times(1)
			},
			want:    entity.Reception{},
			wantErr: service.ErrNoPointFound,
		},
		{
			name: "failed to check point existence",
			mockBehavior: func(r *mock_reception.MockReceptionRepository, t *mock_transactor.MockTransactor, m *mock_reception.MockMetrics) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				r.EXPECT().CheckIfPointExists(ctx, pointID).Return(false, arbitraryErr).Times(1)
				m.EXPECT().ErrInc().Times(1)
			},
			want:    entity.Reception{},
			wantErr: arbitraryErr,
		},
		{
			name: "failed to open",
			mockBehavior: func(r *mock_reception.MockReceptionRepository, t *mock_transactor.MockTransactor, m *mock_reception.MockMetrics) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				r.EXPECT().CheckIfPointExists(ctx, pointID).Return(true, nil).Times(1)
				r.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusClosed, nil).Times(1)
				r.EXPECT().Open(ctx, pointID).Return(entity.Reception{}, arbitraryErr).Times(1)
				m.EXPECT().ErrInc().Times(1)
			},
			want:    entity.Reception{},
			wantErr: arbitraryErr,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			MockReceptionRepo := mock_reception.NewMockReceptionRepository(ctrl)
			MockTransactor := mock_transactor.NewMockTransactor(ctrl)
			MockMetrics := mock_reception.NewMockMetrics(ctrl)

			tc.mockBehavior(MockReceptionRepo, MockTransactor, MockMetrics)

			s := service.New(MockReceptionRepo, MockTransactor, MockMetrics)

			out, err := s.OpenReception(ctx, pointID)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, out)
		})
	}
}

func TestCloseReception(t *testing.T) {
	var (
		ctx                 = context.Background()
		pointID             = uuid.New()
		arbitraryErr        = errors.New("arbitraryErr")
		productsAmount      = 3
		emptyProductsAmount = 0

		emptyStatus entity.ReceptionStatus = ""
	)

	type MockBehavior func(
		r *mock_reception.MockReceptionRepository,
		t *mock_transactor.MockTransactor,
	)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(r *mock_reception.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				r.EXPECT().CheckIfPointExists(ctx, pointID).Return(true, nil).Times(1)
				r.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusInProgress, nil).Times(1)
				r.EXPECT().GetLastReceptionProductsAmount(ctx, pointID).Return(productsAmount, nil)
				r.EXPECT().CloseLastReception(ctx, pointID).Return(nil).Times(1)
			},
			wantErr: nil,
		},
		{
			name: "failed to get status",
			mockBehavior: func(r *mock_reception.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				r.EXPECT().CheckIfPointExists(ctx, pointID).Return(true, nil).Times(1)
				r.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(emptyStatus, arbitraryErr).Times(1)
			},
			wantErr: arbitraryErr,
		},
		{
			name: "last reception is closed",
			mockBehavior: func(r *mock_reception.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				r.EXPECT().CheckIfPointExists(ctx, pointID).Return(true, nil).Times(1)
				r.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusClosed, nil).Times(1)
			},
			wantErr: service.ErrLastReceptionAlreadyClosed,
		},
		{
			name: "cannot close empty reception",
			mockBehavior: func(r *mock_reception.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				r.EXPECT().CheckIfPointExists(ctx, pointID).Return(true, nil).Times(1)
				r.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusInProgress, nil).Times(1)
				r.EXPECT().GetLastReceptionProductsAmount(ctx, pointID).Return(emptyProductsAmount, nil)
			},
			wantErr: service.ErrCannotCloseEmptyReception,
		},
		{
			name: "failed to check products amount",
			mockBehavior: func(r *mock_reception.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				r.EXPECT().CheckIfPointExists(ctx, pointID).Return(true, nil).Times(1)
				r.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusInProgress, nil).Times(1)
				r.EXPECT().GetLastReceptionProductsAmount(ctx, pointID).Return(emptyProductsAmount, arbitraryErr)
			},
			wantErr: arbitraryErr,
		},
		{
			name: "no point found",
			mockBehavior: func(r *mock_reception.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				r.EXPECT().CheckIfPointExists(ctx, pointID).Return(false, nil).Times(1)
			},
			wantErr: service.ErrNoPointFound,
		},
		{
			name: "failed to check point existence",
			mockBehavior: func(r *mock_reception.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				r.EXPECT().CheckIfPointExists(ctx, pointID).Return(false, arbitraryErr).Times(1)
			},
			wantErr: arbitraryErr,
		},
		{
			name: "failed to close",
			mockBehavior: func(r *mock_reception.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				r.EXPECT().CheckIfPointExists(ctx, pointID).Return(true, nil).Times(1)
				r.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusInProgress, nil).Times(1)
				r.EXPECT().GetLastReceptionProductsAmount(ctx, pointID).Return(productsAmount, nil)
				r.EXPECT().CloseLastReception(ctx, pointID).Return(arbitraryErr).Times(1)
			},
			wantErr: arbitraryErr,
		},
		{
			name: "no reception to close found",
			mockBehavior: func(r *mock_reception.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				r.EXPECT().CheckIfPointExists(ctx, pointID).Return(true, nil).Times(1)
				r.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusInProgress, nil).Times(1)
				r.EXPECT().GetLastReceptionProductsAmount(ctx, pointID).Return(productsAmount, nil)
				r.EXPECT().CloseLastReception(ctx, pointID).Return(repository.ErrNoReceptionFound).Times(1)
			},
			wantErr: service.ErrNoReceptionFound,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			MockReceptionRepo := mock_reception.NewMockReceptionRepository(ctrl)
			MockTransactor := mock_transactor.NewMockTransactor(ctrl)
			MockMetrics := mock_reception.NewMockMetrics(ctrl)

			tc.mockBehavior(MockReceptionRepo, MockTransactor)

			s := service.New(MockReceptionRepo, MockTransactor, MockMetrics)

			err := s.CloseReception(ctx, pointID)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
