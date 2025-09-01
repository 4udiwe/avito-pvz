package post_register_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/4udiwe/avito-pvz/internal/api/http/post_register"
	mock_post_register "github.com/4udiwe/avito-pvz/internal/api/http/post_register/mocks"
	"github.com/4udiwe/avito-pvz/internal/auth"
	"github.com/4udiwe/avito-pvz/internal/dto"
	"github.com/4udiwe/avito-pvz/internal/entity"
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
		Role         = entity.RoleModerator
		RefreshToken = "refresh"
		AccessToken  = "access"
		ttl          = 100
	)

	request := post_register.Request{
		Email:    types.Email(Email),
		Password: Password,
		Role:     dto.PostRegisterJSONBodyRole(Role),
	}
	out := auth.Tokens{
		RefreshToken: RefreshToken,
		AccessToken:  AccessToken,
		ExpiresIn:    int64(ttl),
	}
	responseJSON, _ := json.Marshal(out)

	type MockBehavior func(s *mock_post_register.MockUserService)

	for _, tc := range []struct {
		name         string
		mockBehavior MockBehavior
		wantStatus   int
		wantBody     string
	}{
		{
			name: "success",
			mockBehavior: func(s *mock_post_register.MockUserService) {
				s.EXPECT().Register(gomock.Any(), string(request.Email), request.Password, entity.UserRole(request.Role)).Return(&out, nil).Times(1)
			},
			wantStatus: http.StatusCreated,
			wantBody:   string(responseJSON),
		},
		{
			name: "user already exists",
			mockBehavior: func(s *mock_post_register.MockUserService) {
				s.EXPECT().Register(gomock.Any(), string(request.Email), request.Password, entity.UserRole(request.Role)).Return(nil, user.ErrUserAlreadyExists).Times(1)
			},
			wantStatus: http.StatusConflict,
			wantBody:   user.ErrUserAlreadyExists.Error(),
		},
		{
			name: "internal error",
			mockBehavior: func(s *mock_post_register.MockUserService) {
				s.EXPECT().Register(gomock.Any(), string(request.Email), request.Password, entity.UserRole(request.Role)).Return(nil, arbitraryErr).Times(1)
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
			MockService := mock_post_register.NewMockUserService(ctrl)
			tc.mockBehavior(MockService)

			handler := post_register.New(MockService)

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
