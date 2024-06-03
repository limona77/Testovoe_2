package repository

import (
	"context"
	"testovoe_2/internal/model"
	"testovoe_2/pkg/postgres"
)

type IUser interface {
	CreateUser(context.Context, model.User) (model.User, error)
	GetUserByEmail(context.Context, string) (model.User, error)
}

type IToken interface {
	SaveToken(context.Context, model.Token) (model.Token, error)
	GetToken(ctx context.Context, id int) (model.Token, error)
	RemoveToken(ctx context.Context, token string) (int, error)
}
type ISubscribe interface {
	CreateSubscribe(ctx context.Context, subscribe model.Subscribe) (model.Subscribe, error)
	DeleteSubscribe(ctx context.Context, subscribe model.Subscribe) (model.Subscribe, error)
}
type Repositories struct {
	IUser
	IToken
	ISubscribe
}

func NewRepositories(db *postgres.DB) *Repositories {
	return &Repositories{NewUserRepository(db), NewTokenRepository(db), NewSubscribeRepository(db)}
}
