package post_reception_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/4udiwe/avito-pvz/internal/api/http/post_reception"
	mock_post_reception "github.com/4udiwe/avito-pvz/internal/api/http/post_reception/mocks"
	"github.com/4udiwe/avito-pvz/internal/dto"
	"github.com/4udiwe/avito-pvz/internal/entity"
	service "github.com/4udiwe/avito-pvz/internal/service/reception"
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
		arbitraryErr    = errors.New("arbitrary error")
		pointID         = types.UUID(uuid.New())
		receptionID     = types.UUID(uuid.New())
		receptionStatus = entity.ReceptionStatusInProgress
		time            = time.Now()
	)

	request := post_reception.Request{
		PvzId: pointID,
	}
	response := dto.Reception{
		Id:       &receptionID,
		PvzId:    pointID,
		Status:   dto.ReceptionStatus(receptionStatus),
		DateTime: time,
	}
	responseJSON, _ := json.Marshal(response)

	type MockBehavior func(s *mock_post_reception.MockReceptionService)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		wantStatus   int
		wantBody     string
	}{
		{
			name: "success",
			mockBehavior: func(s *mock_post_reception.MockReceptionService) {
				e := entity.Reception{
					ID:        receptionID,
					PointID:   pointID,
					CreatedAt: time,
					Status:    receptionStatus,
				}
				s.EXPECT().OpenReception(gomock.Any(), pointID).Return(e, nil).Times(1)
			},
			wantStatus: http.StatusCreated,
			wantBody:   string(responseJSON),
		},
		{
			name: "no point found",
			mockBehavior: func(s *mock_post_reception.MockReceptionService) {
				s.EXPECT().OpenReception(gomock.Any(), pointID).Return(entity.Reception{}, service.ErrNoPointFound).Times(1)
			},
			wantStatus: http.StatusNotFound,
			wantBody:   service.ErrNoPointFound.Error(),
		},
		{
			name: "last reception not closed",
			mockBehavior: func(s *mock_post_reception.MockReceptionService) {
				s.EXPECT().OpenReception(gomock.Any(), pointID).Return(entity.Reception{}, service.ErrLastReceptionNotClosed).Times(1)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   service.ErrLastReceptionNotClosed.Error(),
		},
		{
			name: "internal error",
			mockBehavior: func(s *mock_post_reception.MockReceptionService) {
				s.EXPECT().OpenReception(gomock.Any(), pointID).Return(entity.Reception{}, arbitraryErr).Times(1)
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
			MockService := mock_post_reception.NewMockReceptionService(ctrl)
			tc.mockBehavior(MockService)

			handler := post_reception.New(MockService)

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
