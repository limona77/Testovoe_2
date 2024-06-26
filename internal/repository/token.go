package repository

import (
	"context"
	"errors"
	"fmt"
	"testovoe_2/internal/custom-errors"
	"testovoe_2/internal/model"
	"testovoe_2/pkg/postgres"

	"github.com/jackc/pgx/v5"
)

type TokenRepository struct {
	*postgres.DB
}

func NewTokenRepository(db *postgres.DB) *TokenRepository {
	return &TokenRepository{db}
}

func (tR *TokenRepository) SaveToken(ctx context.Context, token model.Token) (model.Token, error) {
	path := "internal.repository.token.SaveToken"

	sql, args, err := tR.Builder.
		Insert("public.tokens").
		Columns("user_id", "refresh_token").
		Values(token.UserID, token.RefreshToken).
		Suffix(`
		ON CONFLICT (user_id) DO UPDATE
    SET refresh_token = excluded.refresh_token
		RETURNING id,refresh_token,user_id`).
		ToSql()
	if err != nil {
		return model.Token{}, fmt.Errorf(path+".ToSql, error: {%w}", err)
	}
	var t model.Token
	err = tR.Pool.QueryRow(ctx, sql, args...).
		Scan(&t.ID, &t.RefreshToken, &t.UserID)
	if err != nil {
		return model.Token{}, fmt.Errorf(path+".Scan, error: {%w}", err)
	}

	return t, nil
}

func (tR *TokenRepository) GetToken(ctx context.Context, userId int) (model.Token, error) {
	path := "internal.repository.token.RefreshToken"
	sql, args, err := tR.Builder.Select("id", "refresh_token", "user_id").
		From("public.tokens").
		Where("user_id = ?", userId).
		ToSql()
	if err != nil {
		return model.Token{}, fmt.Errorf(path+".ToSql, error: {%w}", err)
	}

	var modelToken model.Token

	err = tR.Pool.QueryRow(ctx, sql, args...).
		Scan(&modelToken.ID, &modelToken.RefreshToken, &modelToken.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Token{}, custom_errors.ErrUserUnauthorized
		}
		return model.Token{}, fmt.Errorf(path+".QueryRow, error: {%w}", err)
	}
	return modelToken, nil
}

func (tR *TokenRepository) RemoveToken(ctx context.Context, token string) (int, error) {
	path := "internal.repository.token.RemoveToken"
	sql, args, err := tR.Builder.
		Delete("public.tokens").
		Where("refresh_token = ?", token).
		Suffix("RETURNING user_id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf(path+".ToSql, error: {%w}", err)
	}

	var userID int
	err = tR.Pool.QueryRow(ctx, sql, args...).
		Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf(path+".ToSql, error: {%w}", custom_errors.ErrUserUnauthorized)
		}
		return 0, fmt.Errorf(path+".QueryRow, error: {%w}", err)
	}

	return userID, nil
}
