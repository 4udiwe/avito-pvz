package delete_product_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4udiwe/avito-pvz/internal/api/http/delete_product"
	mock_delete_product "github.com/4udiwe/avito-pvz/internal/api/http/delete_product/mocks"
	service "github.com/4udiwe/avito-pvz/internal/service/product"
	"github.com/4udiwe/avito-pvz/pkg/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandle(t *testing.T) {
	var (
		pointID      = uuid.New()
		arbitraryErr = errors.New("arbitrary error")
	)

	type MockBehavior func(s *mock_delete_product.MockProductService)

	for _, tc := range []struct {
		name         string
		pointID      string
		mockBehavior MockBehavior
		wantStatus   int
		wantBody     string
	}{
		{
			name:    "success",
			pointID: pointID.String(),
			mockBehavior: func(s *mock_delete_product.MockProductService) {
				s.EXPECT().DeleteLastProductFromReception(gomock.Any(), pointID).Return(nil).Times(1)
			},
			wantStatus: http.StatusOK,
			wantBody:   "",
		},
		{
			name:    "point not found",
			pointID: pointID.String(),
			mockBehavior: func(s *mock_delete_product.MockProductService) {
				s.EXPECT().DeleteLastProductFromReception(gomock.Any(), pointID).Return(service.ErrNoPointFound).Times(1)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   service.ErrNoPointFound.Error(),
		},
		{
			name:    "reception already closed",
			pointID: pointID.String(),
			mockBehavior: func(s *mock_delete_product.MockProductService) {
				s.EXPECT().DeleteLastProductFromReception(gomock.Any(), pointID).Return(service.ErrReceptionAlreadyClosed).Times(1)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   service.ErrReceptionAlreadyClosed.Error(),
		},
		{
			name:    "no reception found",
			pointID: pointID.String(),
			mockBehavior: func(s *mock_delete_product.MockProductService) {
				s.EXPECT().DeleteLastProductFromReception(gomock.Any(), pointID).Return(service.ErrNoReceptionFound).Times(1)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   service.ErrNoReceptionFound.Error(),
		},
		{
			name:    "internal error",
			pointID: pointID.String(),
			mockBehavior: func(s *mock_delete_product.MockProductService) {
				s.EXPECT().DeleteLastProductFromReception(gomock.Any(), pointID).Return(arbitraryErr).Times(1)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   arbitraryErr.Error(),
		},
		{
			name:         "invalid pvz id provided",
			pointID:      "123",
			mockBehavior: func(s *mock_delete_product.MockProductService) {},
			wantStatus:   http.StatusBadRequest,
			wantBody:      "invalid UUID length: 3",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Validator = validator.NewCustomValidator()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			ctx.SetParamNames("pvzId")
			ctx.SetParamValues(tc.pointID)

			ctrl := gomock.NewController(t)
			MockService := mock_delete_product.NewMockProductService(ctrl)
			tc.mockBehavior(MockService)

			handler := delete_product.New(MockService)

			err := handler.Handle(ctx)

			if tc.wantStatus >= 400 {
				require.Error(t, err)
				httpErr, ok := err.(*echo.HTTPError)
				require.True(t, ok)
				assert.Equal(t, tc.wantStatus, httpErr.Code)
				assert.Equal(t, tc.wantBody, httpErr.Message)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantStatus, rec.Code)
				assert.Equal(t, tc.wantBody, rec.Body.String())
			}
		})
	}
}
