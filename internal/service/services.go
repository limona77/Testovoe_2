package service

import (
	"context"
	"testovoe_2/internal/model"
	"testovoe_2/internal/repository"
)

type AuthParams struct {
	Email    string
	Password string
	Birthday string
}

//go:generate mockgen -source=services.go -destination=mocks/mock.go
type IAuth interface {
	Register(ctx context.Context, params AuthParams) (Tokens, model.User, error)
	GenerateTokens(ctx context.Context, params AuthParams) (Tokens, model.User, error)
	Refresh(ctx context.Context, token string) (Tokens, model.User, error)
	Login(ctx context.Context, params AuthParams) (Tokens, model.User, error)
	Logout(ctx context.Context, token string) (int, error)
	BirthdayChecker(ctx context.Context, params model.User, ch chan []string) error
}
type IClient interface {
	VerifyToken(token string) (TokenClaims, error)
	GetUserByEmail(ctx context.Context, params AuthParams) (model.User, error)
	Subscribe(ctx context.Context, subscribe model.Subscribe) (model.Subscribe, error)
	Unsubscribe(ctx context.Context, subscribe model.Subscribe) error
}

type Services struct {
	IAuth
	IClient
}

type ServicesDeps struct {
	Repository       *repository.Repositories
	SecretKeyAccess  []byte
	SecretKeyRefresh []byte
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{
		IAuth: NewAuthService(
			deps.Repository.IUser,
			deps.Repository.IToken,
			deps.Repository.INotification,
			deps.SecretKeyAccess,
			deps.SecretKeyRefresh),
		IClient: NewClientService(deps.Repository.IUser, deps.Repository.ISubscribe),
	}
}
