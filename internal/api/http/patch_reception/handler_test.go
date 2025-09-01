package patch_reception_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4udiwe/avito-pvz/internal/api/http/patch_reception"
	mock_patch_reception "github.com/4udiwe/avito-pvz/internal/api/http/patch_reception/mocks"
	service "github.com/4udiwe/avito-pvz/internal/service/reception"
	"github.com/4udiwe/avito-pvz/pkg/validator"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandle(t *testing.T) {
	var (
		arbitraryErr = errors.New("arbitrary error")
		pointID      = uuid.New()
	)

	type MockBehavior func(s *mock_patch_reception.MockReceptionService)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		wantStatus   int
		wantBody     string
	}{
		{
			name: "success",
			mockBehavior: func(s *mock_patch_reception.MockReceptionService) {
				s.EXPECT().CloseReception(gomock.Any(), pointID).Return(nil).Times(1)
			},
			wantStatus: http.StatusAccepted,
			wantBody:   "",
		},
		{
			name: "no point found",
			mockBehavior: func(s *mock_patch_reception.MockReceptionService) {
				s.EXPECT().CloseReception(gomock.Any(), pointID).Return(service.ErrNoPointFound).Times(1)
			},
			wantStatus: http.StatusNotFound,
			wantBody:   service.ErrNoPointFound.Error(),
		},
		{
			name: "no reception found",
			mockBehavior: func(s *mock_patch_reception.MockReceptionService) {
				s.EXPECT().CloseReception(gomock.Any(), pointID).Return(service.ErrNoReceptionFound).Times(1)
			},
			wantStatus: http.StatusNotFound,
			wantBody:   service.ErrNoReceptionFound.Error(),
		},
		{
			name: "last reception already closed",
			mockBehavior: func(s *mock_patch_reception.MockReceptionService) {
				s.EXPECT().CloseReception(gomock.Any(), pointID).Return(service.ErrLastReceptionAlreadyClosed).Times(1)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   service.ErrLastReceptionAlreadyClosed.Error(),
		},
		{
			name: "cannot close empty reception",
			mockBehavior: func(s *mock_patch_reception.MockReceptionService) {
				s.EXPECT().CloseReception(gomock.Any(), pointID).Return(service.ErrCannotCloseEmptyReception).Times(1)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   service.ErrCannotCloseEmptyReception.Error(),
		},
		{
			name: "internal error",
			mockBehavior: func(s *mock_patch_reception.MockReceptionService) {
				s.EXPECT().CloseReception(gomock.Any(), pointID).Return(arbitraryErr).Times(1)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   arbitraryErr.Error(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Validator = validator.NewCustomValidator()
			req := httptest.NewRequest(http.MethodPatch, "/", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			ctx.SetParamNames("pvzId")
			ctx.SetParamValues(string(pointID.String()))

			ctrl := gomock.NewController(t)
			MockService := mock_patch_reception.NewMockReceptionService(ctrl)
			tc.mockBehavior(MockService)

			handler := patch_reception.New(MockService)

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
