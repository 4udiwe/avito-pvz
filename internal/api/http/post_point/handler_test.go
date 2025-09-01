package post_point_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/4udiwe/avito-pvz/internal/api/http/post_point"
	mock_post_point "github.com/4udiwe/avito-pvz/internal/api/http/post_point/mocks"
	"github.com/4udiwe/avito-pvz/internal/dto"
	"github.com/4udiwe/avito-pvz/internal/entity"
	service "github.com/4udiwe/avito-pvz/internal/service/point"
	"github.com/4udiwe/avito-pvz/pkg/validator"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandle(t *testing.T) {
	var arbitraryErr = errors.New("arbitrary error")

	id := types.UUID(uuid.New())
	time := time.Now()

	request := post_point.Request{
		City: "Москва",
	}
	response := dto.PVZ{
		Id:               &id,
		City:             "Москва",
		RegistrationDate: &time,
	}
	responseJSON, _ := json.Marshal(response)

	type MockBehavior func(s *mock_post_point.MockPointService)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		wantStatus   int
		wantBody     string
	}{
		{
			name: "success",
			mockBehavior: func(s *mock_post_point.MockPointService) {
				e := entity.Point{
					ID:        *response.Id,
					City:      string(request.City),
					CreatedAt: *response.RegistrationDate,
				}
				s.EXPECT().CreatePoint(gomock.Any(), string(request.City)).Return(e, nil).Times(1)
			},
			wantStatus: http.StatusCreated,
			wantBody:   string(responseJSON),
		},
		{
			name: "no city found",
			mockBehavior: func(s *mock_post_point.MockPointService) {
				s.EXPECT().CreatePoint(gomock.Any(), string(request.City)).Return(entity.Point{}, service.ErrNoCityFound).Times(1)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   service.ErrNoCityFound.Error(),
		},
		{
			name: "internal error",
			mockBehavior: func(s *mock_post_point.MockPointService) {
				s.EXPECT().CreatePoint(gomock.Any(), string(request.City)).Return(entity.Point{}, arbitraryErr).Times(1)
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
			MockService := mock_post_point.NewMockPointService(ctrl)
			tc.mockBehavior(MockService)

			handler := post_point.New(MockService)

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
