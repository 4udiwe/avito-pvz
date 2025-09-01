package post_product_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/4udiwe/avito-pvz/internal/api/http/post_product"
	mock_post_product "github.com/4udiwe/avito-pvz/internal/api/http/post_product/mocks"
	"github.com/4udiwe/avito-pvz/internal/dto"
	"github.com/4udiwe/avito-pvz/internal/entity"
	service "github.com/4udiwe/avito-pvz/internal/service/product"
	"github.com/4udiwe/avito-pvz/pkg/validator"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandle(t *testing.T) {
	var (
		arbitraryErr = errors.New("arbitrary error")
		PvzID        = types.UUID(uuid.New())
		ProductID    = types.UUID(uuid.New())
		ProductType  = entity.ProductTypeElectronics
		ReceptionID  = types.UUID(uuid.New())
		time         = time.Now()

		request = post_product.Request{
			PvzId: PvzID,
			Type:  dto.PostProductsJSONBodyType(ProductType),
		}
		response = dto.Product{
			Id:          &ProductID,
			ReceptionId: ReceptionID,
			DateTime:    &time,
			Type:        dto.ProductType(ProductType),
		}
	)

	responseJSON, _ := json.Marshal(response)

	type MockBehavior func(s *mock_post_product.MockProductService)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		wantStatus   int
		wantBody     string
	}{
		{
			name: "success",
			mockBehavior: func(s *mock_post_product.MockProductService) {
				e := entity.Product{
					ID:          ProductID,
					CreatedAt:   time,
					ReceptionID: ReceptionID,
					Type:        ProductType,
				}
				s.EXPECT().AddProduct(gomock.Any(), request.PvzId, entity.ProductType(request.Type)).Return(e, nil).Times(1)
			},
			wantStatus: http.StatusCreated,
			wantBody:   string(responseJSON),
		},
		{
			name: "no point found",
			mockBehavior: func(s *mock_post_product.MockProductService) {
				s.EXPECT().AddProduct(gomock.Any(), request.PvzId, entity.ProductType(request.Type)).Return(entity.Product{}, service.ErrNoPointFound).Times(1)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   service.ErrNoPointFound.Error(),
		},
		{
			name: "no reception found",
			mockBehavior: func(s *mock_post_product.MockProductService) {
				s.EXPECT().AddProduct(gomock.Any(), request.PvzId, entity.ProductType(request.Type)).Return(entity.Product{}, service.ErrNoReceptionFound).Times(1)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   service.ErrNoReceptionFound.Error(),
		},
		{
			name: "reception already closed",
			mockBehavior: func(s *mock_post_product.MockProductService) {
				s.EXPECT().AddProduct(gomock.Any(), request.PvzId, entity.ProductType(request.Type)).Return(entity.Product{}, service.ErrReceptionAlreadyClosed).Times(1)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   service.ErrReceptionAlreadyClosed.Error(),
		},
		{
			name: "internal error",
			mockBehavior: func(s *mock_post_product.MockProductService) {
				s.EXPECT().AddProduct(gomock.Any(), request.PvzId, entity.ProductType(request.Type)).Return(entity.Product{}, arbitraryErr).Times(1)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   arbitraryErr.Error(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Validator = validator.NewCustomValidator()

			requestBody, _ := json.Marshal(request)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			ctrl := gomock.NewController(t)
			MockService := mock_post_product.NewMockProductService(ctrl)
			tc.mockBehavior(MockService)

			handler := post_product.New(MockService)

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
				assert.Equal(t, tc.wantBody, strings.Trim(rec.Body.String(), "\n"))
			}
		})
	}
}
