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
)

type Repository struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *Repository {
	return &Repository{pg}
}

func (r *Repository) Create(ctx context.Context, user entity.User) (entity.User, error) {
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
				return entity.User{}, repository.ErrUserAlreadyExists
			}
		}
		return entity.User{}, err
	}

	return user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
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
			return entity.User{}, repository.ErrNoUserFound
		}
		return entity.User{}, fmt.Errorf("UserRepository.GetByEmail - Scan: %w", err)
	}
	return user, nil
}
