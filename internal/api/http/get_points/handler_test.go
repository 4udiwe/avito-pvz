package get_points_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4udiwe/avito-pvz/internal/api/http/get_points"
	mock_get_points "github.com/4udiwe/avito-pvz/internal/api/http/get_points/mocks"
	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/4udiwe/avito-pvz/pkg/validator"
	"github.com/go-playground/assert/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandle(t *testing.T) {
	var (
		arbitraryErr = errors.New("arbitrary error")
	)

	type MockBehavior func(s *mock_get_points.MockPointService)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		wantStatus   int
		wantBody     string
	}{
		{
			name: "success",
			mockBehavior: func(s *mock_get_points.MockPointService) {
				s.EXPECT().GetAllPointsFullInfo(gomock.Any()).Return([]entity.PointFullInfo{}, nil).Times(1)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "internal error",
			mockBehavior: func(s *mock_get_points.MockPointService) {
				s.EXPECT().GetAllPointsFullInfo(gomock.Any()).Return(nil, arbitraryErr).Times(1)
			},
			wantStatus: 500,
			wantBody:   arbitraryErr.Error(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Validator = validator.NewCustomValidator()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			ctrl := gomock.NewController(t)
			MockService := mock_get_points.NewMockPointService(ctrl)
			tc.mockBehavior(MockService)

			handler := get_points.New(MockService)

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
			}
		})
	}
}
