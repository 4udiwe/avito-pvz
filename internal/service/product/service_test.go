package product_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4udiwe/avito-pvz/internal/entity"
	mock_transactor "github.com/4udiwe/avito-pvz/internal/mocks"
	"github.com/4udiwe/avito-pvz/internal/repository"
	service "github.com/4udiwe/avito-pvz/internal/service/product"
	"github.com/4udiwe/avito-pvz/internal/service/product/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAddProduct(t *testing.T) {
	var (
		ctx             = context.Background()
		pointID         = uuid.New()
		productType     = entity.ProductTypeElectronics
		lastReceptionID = uuid.New()
		arbitraryErr    = errors.New("arbitraryErr")

		emptyStatus entity.ReceptionStatus = ""
	)

	productOut := entity.Product{
		ID:          uuid.Max,
		ReceptionID: lastReceptionID,
		Type:        productType,
	}

	type MockBehavior func(
		productRepo *mocks.MockProductsRepository,
		receptionRepo *mocks.MockReceptionRepository,
		t *mock_transactor.MockTransactor,
	)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		want         entity.Product
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(productRepo *mocks.MockProductsRepository, receptionRepo *mocks.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				receptionRepo.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusInProgress, nil).Times(1)

				productRepo.EXPECT().Create(ctx, pointID, productType).Return(productOut, nil).Times(1)

			},
			want:    productOut,
			wantErr: nil,
		},
		{
			name: "reception already closed",
			mockBehavior: func(productRepo *mocks.MockProductsRepository, receptionRepo *mocks.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				receptionRepo.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusClosed, nil).Times(1)
			},
			want:    entity.Product{},
			wantErr: service.ErrReceptionAlreadyClosed,
		},
		{
			name: "fetching reception status error",
			mockBehavior: func(productRepo *mocks.MockProductsRepository, receptionRepo *mocks.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				receptionRepo.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(emptyStatus, arbitraryErr).Times(1)
			},
			want:    entity.Product{},
			wantErr: arbitraryErr,
		},
		{
			name: "creating error no reception found",
			mockBehavior: func(productRepo *mocks.MockProductsRepository, receptionRepo *mocks.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				receptionRepo.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusInProgress, nil).Times(1)

				productRepo.EXPECT().Create(ctx, pointID, productType).Return(entity.Product{}, repository.ErrNoReceptionFound).Times(1)
			},
			want:    entity.Product{},
			wantErr: service.ErrNoReceptionFound,
		},
		{
			name: "creating error no point found",
			mockBehavior: func(productRepo *mocks.MockProductsRepository, receptionRepo *mocks.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				receptionRepo.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusInProgress, nil).Times(1)

				productRepo.EXPECT().Create(ctx, pointID, productType).Return(entity.Product{}, repository.ErrNoPointFound).Times(1)
			},
			want:    entity.Product{},
			wantErr: service.ErrNoPointFound,
		},
		{
			name: "creating arbitrary error",
			mockBehavior: func(productRepo *mocks.MockProductsRepository, receptionRepo *mocks.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				receptionRepo.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusInProgress, nil).Times(1)

				productRepo.EXPECT().Create(ctx, pointID, productType).Return(entity.Product{}, arbitraryErr).Times(1)
			},
			want:    entity.Product{},
			wantErr: arbitraryErr,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			MockReceptionRepository := mocks.NewMockReceptionRepository(ctrl)
			MockProductRepository := mocks.NewMockProductsRepository(ctrl)
			MockTransactor := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(MockProductRepository, MockReceptionRepository, MockTransactor)

			s := service.New(MockProductRepository, MockReceptionRepository, MockTransactor)

			out, err := s.AddProduct(ctx, pointID, productType)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, out)
		})
	}

}

func TestDeleteLastProductFromReception(t *testing.T) {
	var (
		ctx          = context.Background()
		pointID      = uuid.New()
		arbitraryErr = errors.New("arbitraryErr")

		emptyStatus entity.ReceptionStatus = ""
	)

	type MockBehavior func(
		productRepo *mocks.MockProductsRepository,
		receptionRepo *mocks.MockReceptionRepository,
		t *mock_transactor.MockTransactor,
	)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(productRepo *mocks.MockProductsRepository, receptionRepo *mocks.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				receptionRepo.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusInProgress, nil).Times(1)

				productRepo.EXPECT().DeleteLastFromReception(ctx, pointID).Return(nil).Times(1)

			},
			wantErr: nil,
		},
		{
			name: "reception already closed",
			mockBehavior: func(productRepo *mocks.MockProductsRepository, receptionRepo *mocks.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				receptionRepo.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusClosed, nil).Times(1)
			},
			wantErr: service.ErrReceptionAlreadyClosed,
		},
		{
			name: "fetching reception status error",
			mockBehavior: func(productRepo *mocks.MockProductsRepository, receptionRepo *mocks.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				receptionRepo.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(emptyStatus, arbitraryErr).Times(1)
			},
			wantErr: arbitraryErr,
		},
		{
			name: "deleting error no reception found",
			mockBehavior: func(productRepo *mocks.MockProductsRepository, receptionRepo *mocks.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				receptionRepo.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusInProgress, nil).Times(1)

				productRepo.EXPECT().DeleteLastFromReception(ctx, pointID).Return(repository.ErrNoReceptionFound).Times(1)
			},
			wantErr: service.ErrNoReceptionFound,
		},
		{
			name: "deleting error no point found",
			mockBehavior: func(productRepo *mocks.MockProductsRepository, receptionRepo *mocks.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				receptionRepo.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusInProgress, nil).Times(1)

				productRepo.EXPECT().DeleteLastFromReception(ctx, pointID).Return(repository.ErrNoPointFound).Times(1)
			},
			wantErr: service.ErrNoPointFound,
		},
		{
			name: "deleting arbitrary error",
			mockBehavior: func(productRepo *mocks.MockProductsRepository, receptionRepo *mocks.MockReceptionRepository, t *mock_transactor.MockTransactor) {
				t.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				receptionRepo.EXPECT().GetLastReceptionStatus(ctx, pointID).Return(entity.ReceptionStatusInProgress, nil).Times(1)

				productRepo.EXPECT().DeleteLastFromReception(ctx, pointID).Return(arbitraryErr).Times(1)
			},
			wantErr: arbitraryErr,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			MockReceptionRepository := mocks.NewMockReceptionRepository(ctrl)
			MockProductRepository := mocks.NewMockProductsRepository(ctrl)
			MockTransactor := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(MockProductRepository, MockReceptionRepository, MockTransactor)

			s := service.New(MockProductRepository, MockReceptionRepository, MockTransactor)

			err := s.DeleteLastProductFromReception(ctx, pointID)
			assert.ErrorIs(t, err, tc.wantErr)
		})

	}
}
