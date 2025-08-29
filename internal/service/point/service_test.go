package point_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/4udiwe/avito-pvz/internal/entity"
	mock_transactor "github.com/4udiwe/avito-pvz/internal/mocks"
	"github.com/4udiwe/avito-pvz/internal/repository"
	service "github.com/4udiwe/avito-pvz/internal/service/point"
	"github.com/4udiwe/avito-pvz/internal/service/point/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreatePoint(t *testing.T) {
	var (
		ctx   = context.Background()
		city  = "Москва"
		point = entity.Point{
			ID:        uuid.Max,
			City:      city,
			CreatedAt: time.Now(),
		}
		emptyPoint   = entity.Point{}
		arbitraryErr = errors.New("arbitrary error")
	)

	type MockBehavior func(r *mocks.MockPointRepository)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		want         entity.Point
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(r *mocks.MockPointRepository) {
				r.EXPECT().Create(ctx, city).Return(point, nil).Times(1)
			},
			want:    point,
			wantErr: nil,
		},
		{
			name: "no city found",
			mockBehavior: func(r *mocks.MockPointRepository) {
				r.EXPECT().Create(ctx, city).Return(emptyPoint, repository.ErrNoCityFound).Times(1)
			},
			want:    emptyPoint,
			wantErr: service.ErrNoCityFound,
		},
		{
			name: "arbitrary error",
			mockBehavior: func(r *mocks.MockPointRepository) {
				r.EXPECT().Create(ctx, city).Return(emptyPoint, arbitraryErr).Times(1)
			},
			want:    emptyPoint,
			wantErr: arbitraryErr,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			MockPointRepository := mocks.NewMockPointRepository(ctrl)
			MockReceptionRepository := mocks.NewMockReceptionRepository(ctrl)
			MockProductRepository := mocks.NewMockProductRepository(ctrl)

			MockTransactor := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(MockPointRepository)

			s := service.New(MockPointRepository, MockReceptionRepository, MockProductRepository, MockTransactor)

			out, err := s.CreatePoint(ctx, city)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, out)
		})
	}
}

func TestGetAllPoints(t *testing.T) {
	var (
		ctx    = context.Background()
		city   = "Москва"
		points = []entity.Point{
			{
				ID:        uuid.New(),
				City:      city,
				CreatedAt: time.Now(),
			},
			{
				ID:        uuid.New(),
				City:      city,
				CreatedAt: time.Now(),
			},
		}
		arbitraryErr = errors.New("arbitrary error")
	)

	type MockBehavior func(r *mocks.MockPointRepository)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		want         []entity.Point
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(r *mocks.MockPointRepository) {
				r.EXPECT().GetAll(ctx).Return(points, nil).Times(1)
			},
			want:    points,
			wantErr: nil,
		},
		{
			name: "cannot fetch points",
			mockBehavior: func(r *mocks.MockPointRepository) {
				r.EXPECT().GetAll(ctx).Return(nil, arbitraryErr).Times(1)
			},
			want:    nil,
			wantErr: arbitraryErr,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			MockPointRepository := mocks.NewMockPointRepository(ctrl)
			MockReceptionRepository := mocks.NewMockReceptionRepository(ctrl)
			MockProductRepository := mocks.NewMockProductRepository(ctrl)

			MockTransactor := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(MockPointRepository)

			s := service.New(MockPointRepository, MockReceptionRepository, MockProductRepository, MockTransactor)

			out, err := s.GetAllPoints(ctx)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, out)
		})
	}
}

func TestGetAllPointsFullInfo(t *testing.T) {
	var (
		ctx          = context.Background()
		arbitraryErr = errors.New("arbitrary error")
	)

	pointID1 := uuid.New()
	pointID2 := uuid.New()
	receptionID1 := uuid.New()
	receptionID2 := uuid.New()

	points := []entity.Point{
		{ID: pointID1, City: "Москва", CreatedAt: time.Now()},
		{ID: pointID2, City: "Санкт-Петербург", CreatedAt: time.Now()},
	}

	receptionsPoint1 := []entity.Reception{
		{ID: receptionID1, PointID: pointID1, Status: entity.ReceptionStatusInProgress, CreatedAt: time.Now()},
	}

	receptionsPoint2 := []entity.Reception{
		{ID: receptionID2, PointID: pointID2, Status: entity.ReceptionStatusClosed, CreatedAt: time.Now()},
	}

	productsReception1 := []entity.Product{
		{ID: uuid.New(), ReceptionID: receptionID1, Type: entity.ProductTypeElectronics, CreatedAt: time.Now()},
		{ID: uuid.New(), ReceptionID: receptionID1, Type: entity.ProductTypeClothes, CreatedAt: time.Now()},
	}

	productsReception2 := []entity.Product{
		{ID: uuid.New(), ReceptionID: receptionID2, Type: entity.ProductTypeShoes, CreatedAt: time.Now()},
	}

	expectedResult := []entity.PointFullInfo{
		{
			Point: points[0],
			Receptions: []entity.ReceptionWithProducts{
				{
					Reception: receptionsPoint1[0],
					Products:  productsReception1,
				},
			},
		},
		{
			Point: points[1],
			Receptions: []entity.ReceptionWithProducts{
				{
					Reception: receptionsPoint2[0],
					Products:  productsReception2,
				},
			},
		},
	}

	type MockBehavior struct {
		pointMock     func(r *mocks.MockPointRepository)
		receptionMock func(r *mocks.MockReceptionRepository)
		productMock   func(r *mocks.MockProductRepository)
	}

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		want         []entity.PointFullInfo
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: MockBehavior{
				pointMock: func(r *mocks.MockPointRepository) {
					r.EXPECT().GetAll(ctx).Return(points, nil).Times(1)
				},
				receptionMock: func(r *mocks.MockReceptionRepository) {
					r.EXPECT().GetAllByPoint(ctx, pointID1).Return(receptionsPoint1, nil).Times(1)
					r.EXPECT().GetAllByPoint(ctx, pointID2).Return(receptionsPoint2, nil).Times(1)
				},
				productMock: func(r *mocks.MockProductRepository) {
					r.EXPECT().GetAllByReception(ctx, receptionID1).Return(productsReception1, nil).Times(1)
					r.EXPECT().GetAllByReception(ctx, receptionID2).Return(productsReception2, nil).Times(1)
				},
			},
			want:    expectedResult,
			wantErr: nil,
		},
		{
			name: "failed to get points",
			mockBehavior: MockBehavior{
				pointMock: func(r *mocks.MockPointRepository) {
					r.EXPECT().GetAll(ctx).Return(nil, arbitraryErr).Times(1)
				},
				receptionMock: func(r *mocks.MockReceptionRepository) {},
				productMock:   func(r *mocks.MockProductRepository) {},
			},
			want:    nil,
			wantErr: arbitraryErr,
		},
		{
			name: "failed to get receptions for point",
			mockBehavior: MockBehavior{
				pointMock: func(r *mocks.MockPointRepository) {
					r.EXPECT().GetAll(ctx).Return(points, nil).Times(1)
				},
				receptionMock: func(r *mocks.MockReceptionRepository) {
					r.EXPECT().GetAllByPoint(ctx, pointID1).Return(nil, arbitraryErr).Times(1)
				},
				productMock: func(r *mocks.MockProductRepository) {},
			},
			want:    nil,
			wantErr: arbitraryErr,
		},
		{
			name: "failed to get products for reception",
			mockBehavior: MockBehavior{
				pointMock: func(r *mocks.MockPointRepository) {
					r.EXPECT().GetAll(ctx).Return(points, nil).Times(1)
				},
				receptionMock: func(r *mocks.MockReceptionRepository) {
					r.EXPECT().GetAllByPoint(ctx, pointID1).Return(receptionsPoint1, nil).Times(1)
				},
				productMock: func(r *mocks.MockProductRepository) {
					r.EXPECT().GetAllByReception(ctx, receptionID1).Return(nil, arbitraryErr).Times(1)
				},
			},
			want:    nil,
			wantErr: arbitraryErr,
		},
		{
			name: "no points found",
			mockBehavior: MockBehavior{
				pointMock: func(r *mocks.MockPointRepository) {
					r.EXPECT().GetAll(ctx).Return([]entity.Point{}, arbitraryErr).Times(1)
				},
				receptionMock: func(r *mocks.MockReceptionRepository) {},
				productMock:   func(r *mocks.MockProductRepository) {},
			},
			want:    nil,
			wantErr: arbitraryErr,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			MockPointRepository := mocks.NewMockPointRepository(ctrl)
			MockReceptionRepository := mocks.NewMockReceptionRepository(ctrl)
			MockProductRepository := mocks.NewMockProductRepository(ctrl)
			MockTransactor := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior.pointMock(MockPointRepository)
			tc.mockBehavior.receptionMock(MockReceptionRepository)
			tc.mockBehavior.productMock(MockProductRepository)

			s := service.New(MockPointRepository, MockReceptionRepository, MockProductRepository, MockTransactor)

			result, err := s.GetAllPointsFullInfo(ctx)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, result)
		})
	}
}
