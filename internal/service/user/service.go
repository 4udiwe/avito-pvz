package user

import (
	"context"
	"errors"

	"github.com/4udiwe/avito-pvz/internal/auth"
	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/4udiwe/avito-pvz/internal/repository"
	"github.com/4udiwe/avito-pvz/pkg/transactor"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Service struct {
	userRepository UserRepository
	txManager      transactor.Transactor
	auth           Auth
	hasher         Hasher
}

func New(r UserRepository, tx transactor.Transactor, a Auth, h Hasher) *Service {
	return &Service{
		userRepository: r,
		txManager:      tx,
		auth:           a,
		hasher:         h,
	}
}

func (s *Service) DummyLogin(ctx context.Context, role entity.UserRole) (string, error) {
	logrus.Infof("Service: generating dummy login with role %s", role)

	user := entity.User{
		Email: "dummyemail@google.com",
		Role:  role,
	}

	tokens, err := s.auth.GenerateTokens(user)
	if err != nil {
		logrus.Errorf("Service: Failed to generate tokens: %v", err)
		return "", err
	}

	logrus.Infof("Service: generated tokens for user dummy user with role %s", role)
	return tokens.RefreshToken, nil
}

func (s *Service) Register(ctx context.Context, email string, password string, role entity.UserRole) (*auth.Tokens, error) {
	logrus.Infof("Service: Registering user %s with role %s", email, role)

	var tokens *auth.Tokens

	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Check if there is no user with this email
		if _, err := s.userRepository.GetByEmail(ctx, email); !errors.Is(err, repository.ErrNoUserFound) {
			return ErrUserAlreadyExists
		}

		// Password hashing
		hash, err := s.hasher.HashPassword(password)

		if err != nil {
			logrus.Errorf("Service: Failed to hash password: %v", err)
			return err
		}

		user := entity.User{
			Email:        email,
			PasswordHash: hash,
			Role:         role,
		}

		// Generate tokens
		tokens, err = s.auth.GenerateTokens(user)
		if err != nil {
			logrus.Errorf("Service: Failed to generate tokens: %v", err)
			return err
		}

		// Assigning refresh token to user and saving to DB
		user.RefreshToken = tokens.RefreshToken

		_, err = s.userRepository.Create(ctx, user)
		if err != nil {
			logrus.Errorf("Service: Failed to create user: %v", err)
			return err
		}
		return err
	})

	if err != nil {
		return nil, err
	}

	logrus.Infof("Service: Registered user %s", email)

	return tokens, nil
}

func (s *Service) Authenticate(ctx context.Context, email string, password string) (*auth.Tokens, error) {
	logrus.Infof("Service: Authenticating user %s", email)

	var tokens *auth.Tokens

	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Receiveing user
		user, err := s.userRepository.GetByEmail(ctx, email)
		if err != nil {
			if errors.Is(err, repository.ErrNoUserFound) {
				logrus.Errorf("Service: No user found by email: %v", err)
				return ErrNoUserFound
			}
			logrus.Errorf("Service: Failed to get user by email: %v", err)
			return err
		}

		// Comparing password hashes
		if !s.hasher.CheckPasswordHash(password, user.PasswordHash) {
			logrus.Warnf("Service: Invalid password for user %s", email)
			return ErrInvalidCredentials
		}

		// Generating new tokens
		tokens, err = s.auth.GenerateTokens(user)
		if err != nil {
			logrus.Errorf("Service: Failed to generate tokens: %v", err)
			return err
		}

		// Updating refresh token for user
		_, err = s.userRepository.UpdateRefreshToken(ctx, user.ID, tokens.RefreshToken)
		if err != nil {
			logrus.Errorf("Service: Failed to update refresh token: %v", err)
		}
		return err
	})

	if err != nil {
		return nil, err
	}

	logrus.Infof("Service: Authenticated user %s", email)
	return tokens, nil
}

func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (*auth.Tokens, error) {
	logrus.Info("Service: Refreshing tokens")

	// Validating refresh token
	email, err := s.auth.ValidateRefreshToken(refreshToken)
	if err != nil {
		logrus.Warnf("Service: Invalid refresh token: %v", err)
		return nil, ErrInvalidRefreshToken
	}
	var tokens *auth.Tokens

	err = s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// Receiving user
		user, err := s.userRepository.GetByEmail(ctx, email)
		if err != nil {
			logrus.Errorf("Service: Failed to get user by email: %v", err)
			return err
		}

		// Check if refresh token is the same as in DB
		if user.RefreshToken != refreshToken {
			logrus.Warnf("Service: Refresh token mismatch for user %s", email)
			return ErrInvalidRefreshToken
		}

		// Generating new tokens
		tokens, err = s.auth.GenerateTokens(user)
		if err != nil {
			logrus.Errorf("Service: Failed to generate tokens: %v", err)
			return err
		}

		// Updating refresh token for user
		_, err = s.userRepository.UpdateRefreshToken(ctx, user.ID, tokens.RefreshToken)
		if err != nil {
			logrus.Errorf("Service: Failed to update refresh token: %v", err)
		}
		return err
	})

	if err != nil {
		return nil, err
	}

	logrus.Infof("Service: Tokens refreshed for user %s", email)
	return tokens, nil
}

func (s *Service) Logout(ctx context.Context, userID uuid.UUID) error {
	logrus.Infof("Service: Logging out user %d", userID)

	// Deleting refresh token (assigning empty string)
	_, err := s.userRepository.UpdateRefreshToken(ctx, userID, "")
	if err != nil {
		logrus.Errorf("Service: Failed to clear refresh token: %v", err)
		return err
	}

	logrus.Infof("Service: User %d logged out", userID)
	return nil
}
