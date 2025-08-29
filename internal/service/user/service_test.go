package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/4udiwe/avito-pvz/internal/auth"
	"github.com/4udiwe/avito-pvz/internal/entity"
	mock_transactor "github.com/4udiwe/avito-pvz/internal/mocks"
	"github.com/4udiwe/avito-pvz/internal/repository"
	service "github.com/4udiwe/avito-pvz/internal/service/user"
	"github.com/4udiwe/avito-pvz/internal/service/user/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRegister(t *testing.T) {
	var (
		ctx          = context.Background()
		arbitraryErr = errors.New("arbitrary error")
		email        = "email123@gmail.com"
		password     = "12345678"
		role         = entity.RoleEmployee
		emptyUser    = entity.User{}
		tokens       = auth.Tokens{
			AccessToken:  "access token",
			RefreshToken: "refresh token",
			ExpiresIn:    100,
		}
		hashedPassword = "hashed_password_123"
	)

	type MockBehavior func(
		u *mocks.MockUserRepository,
		tx *mock_transactor.MockTransactor,
		a *mocks.MockAuth,
		h *mocks.MockHasher,
	)

	for _, tc := range []struct {
		name         string
		email        string
		password     string
		role         entity.UserRole
		mockBehavior MockBehavior
		want         *auth.Tokens
		wantErr      error
	}{
		{
			name:     "success employee",
			email:    email,
			password: password,
			role:     entity.RoleEmployee,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth, h *mocks.MockHasher) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				u.EXPECT().GetByEmail(ctx, email).Return(emptyUser, repository.ErrNoUserFound).Times(1)
				h.EXPECT().HashPassword(password).Return(hashedPassword, nil).Times(1)
				userWithHash := entity.User{
					Email:        email,
					PasswordHash: hashedPassword,
					Role:         entity.RoleEmployee,
				}
				a.EXPECT().GenerateTokens(userWithHash).Return(&tokens, nil).Times(1)
				userToCreate := entity.User{
					Email:        email,
					PasswordHash: hashedPassword,
					Role:         entity.RoleEmployee,
					RefreshToken: tokens.RefreshToken,
				}
				u.EXPECT().Create(ctx, userToCreate).Return(userToCreate, nil).Times(1)
			},
			want:    &tokens,
			wantErr: nil,
		},
		{
			name:     "success moderator",
			email:    "moderator@mail.com",
			password: password,
			role:     entity.RoleModerator,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth, h *mocks.MockHasher) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				u.EXPECT().GetByEmail(ctx, "moderator@mail.com").Return(emptyUser, repository.ErrNoUserFound).Times(1)
				h.EXPECT().HashPassword(password).Return(hashedPassword, nil).Times(1)
				userWithHash := entity.User{
					Email:        "moderator@mail.com",
					PasswordHash: hashedPassword,
					Role:         entity.RoleModerator,
				}
				a.EXPECT().GenerateTokens(userWithHash).Return(&tokens, nil).Times(1)
				userToCreate := entity.User{
					Email:        "moderator@mail.com",
					PasswordHash: hashedPassword,
					Role:         entity.RoleModerator,
					RefreshToken: tokens.RefreshToken,
				}
				u.EXPECT().Create(ctx, userToCreate).Return(userToCreate, nil).Times(1)
			},
			want:    &tokens,
			wantErr: nil,
		},
		{
			name:     "user already exists",
			email:    email,
			password: password,
			role:     role,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth, h *mocks.MockHasher) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				existingUser := entity.User{Email: email, Role: role}
				u.EXPECT().GetByEmail(ctx, email).Return(existingUser, nil).Times(1)
			},
			want:    nil,
			wantErr: service.ErrUserAlreadyExists,
		},
		{
			name:     "hash password error",
			email:    email,
			password: password,
			role:     role,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth, h *mocks.MockHasher) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				u.EXPECT().GetByEmail(ctx, email).Return(emptyUser, repository.ErrNoUserFound).Times(1)
				h.EXPECT().HashPassword(password).Return("", arbitraryErr).Times(1)
			},
			want:    nil,
			wantErr: arbitraryErr,
		},
		{
			name:     "generate tokens error",
			email:    email,
			password: password,
			role:     role,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth, h *mocks.MockHasher) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				u.EXPECT().GetByEmail(ctx, email).Return(emptyUser, repository.ErrNoUserFound).Times(1)
				h.EXPECT().HashPassword(password).Return(hashedPassword, nil).Times(1)
				userWithHash := entity.User{
					Email:        email,
					PasswordHash: hashedPassword,
					Role:         role,
				}
				a.EXPECT().GenerateTokens(userWithHash).Return(nil, arbitraryErr).Times(1)
			},
			want:    nil,
			wantErr: arbitraryErr,
		},
		{
			name:     "create user error",
			email:    email,
			password: password,
			role:     role,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth, h *mocks.MockHasher) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				u.EXPECT().GetByEmail(ctx, email).Return(emptyUser, repository.ErrNoUserFound).Times(1)
				h.EXPECT().HashPassword(password).Return(hashedPassword, nil).Times(1)
				userWithHash := entity.User{
					Email:        email,
					PasswordHash: hashedPassword,
					Role:         role,
				}
				a.EXPECT().GenerateTokens(userWithHash).Return(&tokens, nil).Times(1)
				userToCreate := entity.User{
					Email:        email,
					PasswordHash: hashedPassword,
					Role:         role,
					RefreshToken: tokens.RefreshToken,
				}
				u.EXPECT().Create(ctx, userToCreate).Return(emptyUser, arbitraryErr).Times(1)
			},
			want:    nil,
			wantErr: arbitraryErr,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			MockUserRepo := mocks.NewMockUserRepository(ctrl)
			MockAuth := mocks.NewMockAuth(ctrl)
			MockHasher := mocks.NewMockHasher(ctrl)
			MockTransactor := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(MockUserRepo, MockTransactor, MockAuth, MockHasher)

			s := service.New(MockUserRepo, MockTransactor, MockAuth, MockHasher)

			out, err := s.Register(ctx, tc.email, tc.password, tc.role)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, out)
		})
	}
}

func TestAuthenticate(t *testing.T) {
	var (
		ctx            = context.Background()
		arbitraryErr   = errors.New("arbitrary error")
		email          = "email123@gmail.com"
		password       = "12345678"
		hashedPassword = "$2a$10$hashedpassword123"
		userID         = uuid.New()
	)

	validUser := entity.User{
		ID:           userID,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         entity.RoleEmployee,
	}

	tokens := &auth.Tokens{
		AccessToken:  "access_token_123",
		RefreshToken: "refresh_token_456",
		ExpiresIn:    3600,
	}

	type MockBehavior func(
		u *mocks.MockUserRepository,
		tx *mock_transactor.MockTransactor,
		a *mocks.MockAuth,
		h *mocks.MockHasher,
	)

	for _, tc := range []struct {
		name         string
		email        string
		password     string
		mockBehavior MockBehavior
		want         *auth.Tokens
		wantErr      error
	}{
		{
			name:     "success employee",
			email:    email,
			password: password,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth, h *mocks.MockHasher) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				u.EXPECT().GetByEmail(ctx, email).Return(validUser, nil).Times(1)
				h.EXPECT().CheckPasswordHash(password, hashedPassword).Return(true).Times(1)
				a.EXPECT().GenerateTokens(validUser).Return(tokens, nil).Times(1)
				u.EXPECT().UpdateRefreshToken(ctx, userID, tokens.RefreshToken).Return(validUser, nil).Times(1)
			},
			want:    tokens,
			wantErr: nil,
		},
		{
			name:     "success moderator",
			email:    "moderator@mail.com",
			password: password,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth, h *mocks.MockHasher) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				moderatorUser := entity.User{
					ID:           uuid.New(),
					Email:        "moderator@mail.com",
					PasswordHash: hashedPassword,
					Role:         entity.RoleModerator,
				}
				u.EXPECT().GetByEmail(ctx, "moderator@mail.com").Return(moderatorUser, nil).Times(1)
				h.EXPECT().CheckPasswordHash(password, hashedPassword).Return(true).Times(1)
				a.EXPECT().GenerateTokens(moderatorUser).Return(tokens, nil).Times(1)
				u.EXPECT().UpdateRefreshToken(ctx, moderatorUser.ID, tokens.RefreshToken).Return(moderatorUser, nil).Times(1)
			},
			want:    tokens,
			wantErr: nil,
		},
		{
			name:     "user not found",
			email:    email,
			password: password,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth, h *mocks.MockHasher) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				u.EXPECT().GetByEmail(ctx, email).Return(entity.User{}, repository.ErrNoUserFound).Times(1)
			},
			want:    nil,
			wantErr: service.ErrNoUserFound,
		},
		{
			name:     "get user error",
			email:    email,
			password: password,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth, h *mocks.MockHasher) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				u.EXPECT().GetByEmail(ctx, email).Return(entity.User{}, arbitraryErr).Times(1)
			},
			want:    nil,
			wantErr: arbitraryErr,
		},
		{
			name:     "invalid password",
			email:    email,
			password: "wrong_password",
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth, h *mocks.MockHasher) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				u.EXPECT().GetByEmail(ctx, email).Return(validUser, nil).Times(1)
				h.EXPECT().CheckPasswordHash("wrong_password", hashedPassword).Return(false).Times(1)
			},
			want:    nil,
			wantErr: service.ErrInvalidCredentials,
		},
		{
			name:     "generate tokens error",
			email:    email,
			password: password,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth, h *mocks.MockHasher) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				u.EXPECT().GetByEmail(ctx, email).Return(validUser, nil).Times(1)
				h.EXPECT().CheckPasswordHash(password, hashedPassword).Return(true).Times(1)
				a.EXPECT().GenerateTokens(validUser).Return(nil, arbitraryErr).Times(1)
			},
			want:    nil,
			wantErr: arbitraryErr,
		},
		{
			name:     "update refresh token error",
			email:    email,
			password: password,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth, h *mocks.MockHasher) {
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				u.EXPECT().GetByEmail(ctx, email).Return(validUser, nil).Times(1)
				h.EXPECT().CheckPasswordHash(password, hashedPassword).Return(true).Times(1)
				a.EXPECT().GenerateTokens(validUser).Return(tokens, nil).Times(1)
				u.EXPECT().UpdateRefreshToken(ctx, userID, tokens.RefreshToken).Return(entity.User{}, arbitraryErr).Times(1)
			},
			want:    nil,
			wantErr: arbitraryErr,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			MockUserRepo := mocks.NewMockUserRepository(ctrl)
			MockAuth := mocks.NewMockAuth(ctrl)
			MockHasher := mocks.NewMockHasher(ctrl)
			MockTransactor := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(MockUserRepo, MockTransactor, MockAuth, MockHasher)

			s := service.New(MockUserRepo, MockTransactor, MockAuth, MockHasher)

			out, err := s.Authenticate(ctx, tc.email, tc.password)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, out)
		})
	}
}

func TestRefreshTokens(t *testing.T) {
	var (
		ctx             = context.Background()
		arbitraryErr    = errors.New("arbitrary error")
		email           = "email123@gmail.com"
		refreshToken    = "valid_refresh_token_123"
		oldRefreshToken = "old_refresh_token_456"
		userID          = uuid.New()
	)

	validUser := entity.User{
		ID:           userID,
		Email:        email,
		Role:         entity.RoleEmployee,
		RefreshToken: refreshToken,
	}

	userWithDifferentToken := entity.User{
		ID:           userID,
		Email:        email,
		Role:         entity.RoleEmployee,
		RefreshToken: oldRefreshToken, // отличается от переданного токена
	}

	tokens := &auth.Tokens{
		AccessToken:  "new_access_token_123",
		RefreshToken: "new_refresh_token_456",
		ExpiresIn:    3600,
	}

	type MockBehavior func(
		u *mocks.MockUserRepository,
		tx *mock_transactor.MockTransactor,
		a *mocks.MockAuth,
	)

	for _, tc := range []struct {
		name         string
		refreshToken string
		mockBehavior MockBehavior
		want         *auth.Tokens
		wantErr      error
	}{
		{
			name:         "success",
			refreshToken: refreshToken,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth) {
				a.EXPECT().ValidateRefreshToken(refreshToken).Return(email, nil).Times(1)
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				u.EXPECT().GetByEmail(ctx, email).Return(validUser, nil).Times(1)
				a.EXPECT().GenerateTokens(validUser).Return(tokens, nil).Times(1)
				u.EXPECT().UpdateRefreshToken(ctx, userID, tokens.RefreshToken).Return(validUser, nil).Times(1)
			},
			want:    tokens,
			wantErr: nil,
		},
		{
			name:         "success different role",
			refreshToken: refreshToken,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth) {
				a.EXPECT().ValidateRefreshToken(refreshToken).Return("moderator@mail.com", nil).Times(1)
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				moderatorUser := entity.User{
					ID:           uuid.New(),
					Email:        "moderator@mail.com",
					Role:         entity.RoleModerator,
					RefreshToken: refreshToken,
				}
				u.EXPECT().GetByEmail(ctx, "moderator@mail.com").Return(moderatorUser, nil).Times(1)
				a.EXPECT().GenerateTokens(moderatorUser).Return(tokens, nil).Times(1)
				u.EXPECT().UpdateRefreshToken(ctx, moderatorUser.ID, tokens.RefreshToken).Return(moderatorUser, nil).Times(1)
			},
			want:    tokens,
			wantErr: nil,
		},
		{
			name:         "invalid refresh token",
			refreshToken: "invalid_token",
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth) {
				a.EXPECT().ValidateRefreshToken("invalid_token").Return("", arbitraryErr).Times(1)
				// Транзакция не должна начинаться
			},
			want:    nil,
			wantErr: service.ErrInvalidRefreshToken,
		},
		{
			name:         "user not found",
			refreshToken: refreshToken,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth) {
				a.EXPECT().ValidateRefreshToken(refreshToken).Return(email, nil).Times(1)
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				u.EXPECT().GetByEmail(ctx, email).Return(entity.User{}, repository.ErrNoUserFound).Times(1)
				// GenerateTokens и UpdateRefreshToken не должны вызываться
			},
			want:    nil,
			wantErr: repository.ErrNoUserFound,
		},
		{
			name:         "get user error",
			refreshToken: refreshToken,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth) {
				a.EXPECT().ValidateRefreshToken(refreshToken).Return(email, nil).Times(1)
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				u.EXPECT().GetByEmail(ctx, email).Return(entity.User{}, arbitraryErr).Times(1)
				// GenerateTokens и UpdateRefreshToken не должны вызываться
			},
			want:    nil,
			wantErr: arbitraryErr,
		},
		{
			name:         "refresh token mismatch",
			refreshToken: refreshToken,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth) {
				a.EXPECT().ValidateRefreshToken(refreshToken).Return(email, nil).Times(1)
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				u.EXPECT().GetByEmail(ctx, email).Return(userWithDifferentToken, nil).Times(1)
				// GenerateTokens и UpdateRefreshToken не должны вызываться
			},
			want:    nil,
			wantErr: service.ErrInvalidRefreshToken,
		},
		{
			name:         "generate tokens error",
			refreshToken: refreshToken,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth) {
				a.EXPECT().ValidateRefreshToken(refreshToken).Return(email, nil).Times(1)
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				u.EXPECT().GetByEmail(ctx, email).Return(validUser, nil).Times(1)
				a.EXPECT().GenerateTokens(validUser).Return(nil, arbitraryErr).Times(1)
				// UpdateRefreshToken не должен вызываться
			},
			want:    nil,
			wantErr: arbitraryErr,
		},
		{
			name:         "update refresh token error",
			refreshToken: refreshToken,
			mockBehavior: func(u *mocks.MockUserRepository, tx *mock_transactor.MockTransactor, a *mocks.MockAuth) {
				a.EXPECT().ValidateRefreshToken(refreshToken).Return(email, nil).Times(1)
				tx.EXPECT().WithinTransaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				u.EXPECT().GetByEmail(ctx, email).Return(validUser, nil).Times(1)
				a.EXPECT().GenerateTokens(validUser).Return(tokens, nil).Times(1)
				u.EXPECT().UpdateRefreshToken(ctx, userID, tokens.RefreshToken).Return(entity.User{}, arbitraryErr).Times(1)
			},
			want:    nil,
			wantErr: arbitraryErr,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			MockUserRepo := mocks.NewMockUserRepository(ctrl)
			MockAuth := mocks.NewMockAuth(ctrl)
			MockHasher := mocks.NewMockHasher(ctrl)
			MockTransactor := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(MockUserRepo, MockTransactor, MockAuth)

			s := service.New(MockUserRepo, MockTransactor, MockAuth, MockHasher)

			out, err := s.RefreshTokens(ctx, tc.refreshToken)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, out)
		})
	}
}

func TestDummyLogin(t *testing.T) {
	var (
		ctx          = context.Background()
		arbitraryErr = errors.New("arbitrary error")
	)

	tokens := &auth.Tokens{
		AccessToken:  "dummy_access_token",
		RefreshToken: "dummy_refresh_token",
		ExpiresIn:    3600,
	}

	type MockBehavior func(a *mocks.MockAuth)

	for _, tc := range []struct {
		name         string
		role         entity.UserRole
		mockBehavior MockBehavior
		want         string
		wantErr      error
	}{
		{
			name: "success employee role",
			role: entity.RoleEmployee,
			mockBehavior: func(a *mocks.MockAuth) {
				expectedUser := entity.User{
					Email: "dummyemail@google.com",
					Role:  entity.RoleEmployee,
				}
				a.EXPECT().GenerateTokens(expectedUser).Return(tokens, nil).Times(1)
			},
			want:    tokens.RefreshToken,
			wantErr: nil,
		},
		{
			name: "success moderator role",
			role: entity.RoleModerator,
			mockBehavior: func(a *mocks.MockAuth) {
				expectedUser := entity.User{
					Email: "dummyemail@google.com",
					Role:  entity.RoleModerator,
				}
				a.EXPECT().GenerateTokens(expectedUser).Return(tokens, nil).Times(1)
			},
			want:    tokens.RefreshToken,
			wantErr: nil,
		},
		{
			name: "generate tokens error",
			role: entity.RoleEmployee,
			mockBehavior: func(a *mocks.MockAuth) {
				expectedUser := entity.User{
					Email: "dummyemail@google.com",
					Role:  entity.RoleEmployee,
				}
				a.EXPECT().GenerateTokens(expectedUser).Return(nil, arbitraryErr).Times(1)
			},
			want:    "",
			wantErr: arbitraryErr,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			MockUserRepo := mocks.NewMockUserRepository(ctrl)
			MockAuth := mocks.NewMockAuth(ctrl)
			MockHasher := mocks.NewMockHasher(ctrl)
			MockTransactor := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(MockAuth)

			s := service.New(MockUserRepo, MockTransactor, MockAuth, MockHasher)

			out, err := s.DummyLogin(ctx, tc.role)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, out)
		})
	}
}

func TestLogout(t *testing.T) {
	var (
		ctx          = context.Background()
		userID       = uuid.New()
		arbitraryErr = errors.New("arbitrary error")
	)

	updatedUser := entity.User{
		ID:           userID,
		Email:        "test@mail.com",
		Role:         entity.RoleEmployee,
		RefreshToken: "", // после логаута должен быть nil или пустая строка
	}

	type MockBehavior func(u *mocks.MockUserRepository)

	for _, tc := range []struct {
		name         string
		userID       uuid.UUID
		mockBehavior MockBehavior
		wantErr      error
	}{
		{
			name:   "success",
			userID: userID,
			mockBehavior: func(u *mocks.MockUserRepository) {
				u.EXPECT().UpdateRefreshToken(ctx, userID, "").Return(updatedUser, nil).Times(1)
			},
			wantErr: nil,
		},
		{
			name:   "update refresh token error",
			userID: userID,
			mockBehavior: func(u *mocks.MockUserRepository) {
				u.EXPECT().UpdateRefreshToken(ctx, userID, "").Return(entity.User{}, arbitraryErr).Times(1)
			},
			wantErr: arbitraryErr,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			MockUserRepo := mocks.NewMockUserRepository(ctrl)
			MockAuth := mocks.NewMockAuth(ctrl)
			MockHasher := mocks.NewMockHasher(ctrl)
			MockTransactor := mock_transactor.NewMockTransactor(ctrl)

			tc.mockBehavior(MockUserRepo)

			s := service.New(MockUserRepo, MockTransactor, MockAuth, MockHasher)

			err := s.Logout(ctx, tc.userID)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
