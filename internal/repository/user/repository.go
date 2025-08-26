package repo_user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/4udiwe/avito-pvz/internal/repository"
	"github.com/4udiwe/avito-pvz/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *Repository {
	return &Repository{pg}
}

func (r *Repository) Create(ctx context.Context, user entity.User) (entity.User, error) {
	logrus.Infof("Attempting to create user: %s", user.Email)

	query, args, _ := r.Builder.Insert("users").
		Columns("email", "password_hash", "role", "refresh_token").
		Values(user.Email, user.PasswordHash, user.Role, user.RefreshToken).
		Suffix("RETURNING id").
		ToSql()

	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(&user.ID)

	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				logrus.Warnf("User already exists: %s", user.Email)
				return entity.User{}, repository.ErrUserAlreadyExists
			}
		}
		logrus.Errorf("Failed to create user %s: %v", user.Email, err)
		return entity.User{}, err
	}

	logrus.Infof("User created: %+v", user)
	return user, nil
}

func (r *Repository) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) (entity.User, error) {
	logrus.Infof("Updating refresh token for user %d", userID)

	query, args, _ := r.Builder.Update("users").
		Set("refresh_token", refreshToken).
		Set("updated_at", time.Now()).
		Where("id = ?", userID).
		Suffix("RETURNING id, email, role, created_at, updated_at").
		ToSql()

	var user entity.User
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		logrus.Errorf("Failed to update refresh token for user %d: %v", userID, err)
		return entity.User{}, err
	}

	logrus.Infof("Refresh token updated for user %d", userID)
	return user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	logrus.Infof("Fetching user by email: %s", email)

	query, args, _ := r.Builder.
		Select("id", "password_hash", "role", "refresh_token", "created_at", "updated_at").
		From("users").
		Where("email = ?", email).
		ToSql()

	user := entity.User{Email: email}
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.PasswordHash,
		&user.Role,
		&user.RefreshToken,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logrus.Warnf("No user found with email: %s", email)
			return entity.User{}, repository.ErrNoUserFound
		}
		logrus.Errorf("Failed to fetch user %s: %v", email, err)
		return entity.User{}, fmt.Errorf("UserRepository.GetByEmail - Scan: %w", err)
	}

	logrus.Infof("Fetched user: %+v", user)
	return user, nil
}
