package post_login_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/4udiwe/avito-pvz/internal/api/http/post_login"
	mock_post_login "github.com/4udiwe/avito-pvz/internal/api/http/post_login/mocks"
	"github.com/4udiwe/avito-pvz/internal/auth"
	"github.com/4udiwe/avito-pvz/internal/service/user"
	"github.com/4udiwe/avito-pvz/pkg/validator"
	"github.com/go-playground/assert/v2"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandle(t *testing.T) {
	var (
		arbitraryErr = errors.New("arbitrary error")
		Email        = "example@gmail.gom"
		Password     = "12345678"
		RefreshToken = "refresh"
		AccessToken  = "access"
		ttl          = 100
	)

	request := post_login.Request{
		Email:    types.Email(Email),
		Password: Password,
	}
	out := auth.Tokens{
		RefreshToken: RefreshToken,
		AccessToken:  AccessToken,
		ExpiresIn:    int64(ttl),
	}
	responseJSON, _ := json.Marshal(out)

	type MockBehavior func(s *mock_post_login.MockUserService)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		wantStatus   int
		wantBody     string
	}{
		{
			name: "success",
			mockBehavior: func(s *mock_post_login.MockUserService) {
				s.EXPECT().Authenticate(gomock.Any(), string(request.Email), request.Password).Return(&out, nil).Times(1)
			},
			wantStatus: http.StatusCreated,
			wantBody:   string(responseJSON),
		},
		{
			name: "no user found",
			mockBehavior: func(s *mock_post_login.MockUserService) {
				s.EXPECT().Authenticate(gomock.Any(), string(request.Email), request.Password).Return(nil, user.ErrNoUserFound).Times(1)
			},
			wantStatus: http.StatusNotFound,
			wantBody:   user.ErrNoUserFound.Error(),
		},
		{
			name: "invalid credentials",
			mockBehavior: func(s *mock_post_login.MockUserService) {
				s.EXPECT().Authenticate(gomock.Any(), string(request.Email), request.Password).Return(nil, user.ErrInvalidCredentials).Times(1)
			},
			wantStatus: http.StatusForbidden,
			wantBody:   user.ErrInvalidCredentials.Error(),
		},
		{
			name: "internal error",
			mockBehavior: func(s *mock_post_login.MockUserService) {
				s.EXPECT().Authenticate(gomock.Any(), string(request.Email), request.Password).Return(nil, arbitraryErr).Times(1)
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
			MockService := mock_post_login.NewMockUserService(ctrl)
			tc.mockBehavior(MockService)

			handler := post_login.New(MockService)

			err := handler.Handle(ctx)

			if tc.wantStatus >= 400 {
				require.Error(t, err)
				httpErr := &echo.HTTPError{}
				ok := errors.As(err, &httpErr)
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
