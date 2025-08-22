package repo_user

import (
	"context"
	"errors"
	"fmt"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/4udiwe/avito-pvz/internal/repository"
	"github.com/4udiwe/avito-pvz/pkg/postgres"
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
		Columns("email", "password", "role").
		Values(user.Email, user.Password, user.Role).
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

func (r *Repository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	logrus.Infof("Fetching user by email: %s", email)

	query, args, _ := r.Builder.
		Select("id", "password", "role").
		From("users").
		Where("email = ?", email).
		ToSql()

	user := entity.User{Email: email}
	err := r.GetTxManager(ctx).QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Password,
		&user.Role,
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
