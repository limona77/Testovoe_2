package repository

import (
	"context"
	"errors"
	"fmt"
	custom_errros "testovoe_2/internal/custom-errors"
	"testovoe_2/internal/model"
	"testovoe_2/pkg/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository struct {
	*postgres.DB
}

func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{db}
}

func (uR *UserRepository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	path := "internal.repository.user.CreateUser"
	sql, args, err := uR.Builder.
		Insert("public.users").
		Columns("email", "password", "birthday").
		Values(user.Email, user.Password, user.Birthday).
		Suffix("RETURNING id, email, birthday").
		ToSql()
	if err != nil {
		return model.User{}, fmt.Errorf(path+".ToSql, error: {%w}", err)
	}
	var u model.User
	err = uR.Pool.QueryRow(ctx, sql, args...).
		Scan(&u.ID, &u.Email, &u.Birthday)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return model.User{}, custom_errros.ErrAlreadyExists
			}
		}
		return u, fmt.Errorf(path+".QueryRow, error: {%w}", err)
	}
	return u, nil
}

func (uR *UserRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	path := "internal.repository.user.GetUserByEmail"
	sql, args, err := uR.Builder.
		Select("id,email,password").
		From("public.users").
		Where("email = ?", email).
		ToSql()
	if err != nil {
		return model.User{}, fmt.Errorf(path+".ToSql, error: {%w}", err)
	}

	var user model.User
	err = uR.Pool.QueryRow(ctx, sql, args...).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			return model.User{}, err
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, custom_errros.ErrUserNotFound
		}
		return model.User{}, fmt.Errorf(path+".QueryRow, error: {%w}", err)
	}

	return user, nil
}
