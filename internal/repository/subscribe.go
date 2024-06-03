package repository

import (
	"context"
	"errors"
	"fmt"
	custom_errros "testovoe_2/internal/custom-errors"
	"testovoe_2/internal/model"
	"testovoe_2/pkg/postgres"

	"github.com/jackc/pgx/v5/pgconn"
)

type SubscribeRepository struct {
	*postgres.DB
}

func NewSubscribeRepository(db *postgres.DB) *SubscribeRepository {
	return &SubscribeRepository{db}
}

func (sR *SubscribeRepository) CreateSubscribe(
	ctx context.Context,
	subscribe model.Subscribe,
) (model.Subscribe, error) {
	path := "internal.repository.subscribe.CreateSubscribe"
	sql, args, err := sR.Builder.
		Insert("public.subscriptions").
		Columns("user_id", "subscribed_to_id").
		Values(subscribe.UserID, subscribe.SubscribedToId).
		Suffix("RETURNING user_id, subscribed_to_id").
		ToSql()
	if err != nil {
		return model.Subscribe{}, fmt.Errorf(path+".ToSql, error: {%w}", err)
	}
	var s model.Subscribe
	err = sR.Pool.QueryRow(ctx, sql, args...).
		Scan(&s.UserID, &s.SubscribedToId)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return model.Subscribe{}, custom_errros.ErrYouAlreadySub
			}
		}
		return s, fmt.Errorf(path+".QueryRow, error: {%w}", err)
	}
	return s, nil
}

func (sR *SubscribeRepository) DeleteSubscribe(
	ctx context.Context,
	subscribe model.Subscribe,
) (model.Subscribe, error) {
	return model.Subscribe{}, nil
}
