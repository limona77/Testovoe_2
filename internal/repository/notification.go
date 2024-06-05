package repository

import (
	"context"
	"errors"
	"fmt"
	custom_errros "testovoe_2/internal/custom-errors"
	"testovoe_2/internal/model"
	"testovoe_2/pkg/postgres"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type NotificationRepository struct {
	*postgres.DB
}

func NewNotificationRepository(db *postgres.DB) *NotificationRepository {
	return &NotificationRepository{db}
}

func (nR *NotificationRepository) CheckBirthdays(ctx context.Context, params model.User) ([]string, error) {
	path := "internal.repository.notification.CheckBirthdays"
	birthday := time.Now()
	sql, args, err := nR.Builder.
		Select("users.birthday, email").
		From("public.users").
		Join("public.subscriptions ON users.id = subscriptions.subscribed_to_id").
		Where("subscriptions.user_id = ?", params.ID).
		ToSql()
	if err != nil {
		return []string{}, fmt.Errorf(path+".ToSql, error: {%w}", err)
	}

	rows, err := nR.Pool.Query(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			return []string{}, err
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return []string{}, custom_errros.ErrUserNotFound
		}
		return []string{}, fmt.Errorf(path+".QueryRow, error: {%w}", err)
	}
	res := make([]string, 0)
	for rows.Next() {
		b := struct {
			Birthday time.Time
			Email    string
		}{}
		err := rows.Scan(&b.Birthday, &b.Email)
		if err != nil {
			return []string{}, fmt.Errorf(path+".Scan, error: {%w}", err)
		}
		if b.Birthday.Day() == birthday.Day() && b.Birthday.Month() == birthday.Month() {
			res = append(res, fmt.Sprintf("ДР у пользователя %s: %s", b.Email, b.Birthday.Format("02-01-2006")))
		}
	}
	fmt.Println(res)
	return res, nil
}
