package post_refresh_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/4udiwe/avito-pvz/internal/api/http/post_refresh"
	mock_post_refresh "github.com/4udiwe/avito-pvz/internal/api/http/post_refresh/mocks"
	"github.com/4udiwe/avito-pvz/internal/auth"
	"github.com/4udiwe/avito-pvz/internal/service/user"
	"github.com/4udiwe/avito-pvz/pkg/validator"
	"github.com/go-playground/assert/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandle(t *testing.T) {
	var (
		arbitraryErr    = errors.New("arbitrary error")
		OldRefreshToken = "old_refresh"
		RefreshToken    = "refresh"
		AccessToken     = "access"
		ttl             = 100
	)

	request := post_refresh.Request{
		RefreshToken: OldRefreshToken,
	}
	out := auth.Tokens{
		RefreshToken: RefreshToken,
		AccessToken:  AccessToken,
		ExpiresIn:    int64(ttl),
	}
	responseJSON, _ := json.Marshal(out)

	type MockBehavior func(s *mock_post_refresh.MockUserService)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		wantStatus   int
		wantBody     string
	}{
		{
			name: "success",
			mockBehavior: func(s *mock_post_refresh.MockUserService) {
				s.EXPECT().RefreshTokens(gomock.Any(), request.RefreshToken).Return(&out, nil).Times(1)
			},
			wantStatus: http.StatusCreated,
			wantBody:   string(responseJSON),
		},
		{
			name: "invalid refresh token",
			mockBehavior: func(s *mock_post_refresh.MockUserService) {
				s.EXPECT().RefreshTokens(gomock.Any(), request.RefreshToken).Return(&auth.Tokens{}, user.ErrInvalidRefreshToken).Times(1)
			},
			wantStatus: http.StatusForbidden,
			wantBody:   user.ErrInvalidRefreshToken.Error(),
		},
		{
			name: "internal error",
			mockBehavior: func(s *mock_post_refresh.MockUserService) {
				s.EXPECT().RefreshTokens(gomock.Any(), request.RefreshToken).Return(&auth.Tokens{}, arbitraryErr).Times(1)
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
			MockService := mock_post_refresh.NewMockUserService(ctrl)
			tc.mockBehavior(MockService)

			handler := post_refresh.New(MockService)

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
