package service

import (
	"context"
	"fmt"
	custom_errors "testovoe_2/internal/custom-errors"
	"testovoe_2/internal/model"
	"testovoe_2/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ClientService struct {
	userRepository      repository.IUser
	subscribeRepository repository.ISubscribe
}

func NewClientService(uR repository.IUser, sR repository.ISubscribe) *ClientService {
	return &ClientService{uR, sR}
}

func (cS *ClientService) VerifyToken(token string) (TokenClaims, error) {
	path := "internal.service.auth.ParseToken"

	var tokenClaims TokenClaims

	t, err := jwt.ParseWithClaims(token, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
		return tokenClaims.Key, nil
	})
	if err != nil {
		return TokenClaims{}, fmt.Errorf(path+".ParseWithClaims, error: {%w}", err)
	}

	if !t.Valid {
		return tokenClaims, fmt.Errorf(path+".Valid, error: {%w}", err)
	}
	if time.Now().Unix() > tokenClaims.Exp {
		return tokenClaims, custom_errors.ErrTokenExpired
	}

	return tokenClaims, nil
}

func (cS *ClientService) GetUserByEmail(ctx context.Context, params AuthParams) (model.User, error) {
	path := "internal.service.auth.GetUserByEmail"
	user, err := cS.userRepository.GetUserByEmail(ctx, params.Email)
	if err != nil {
		return model.User{}, fmt.Errorf(path+".GetUserByEmail, error: {%w}", err)
	}
	user.Password = ""
	return user, nil
}

func (cS *ClientService) Subscribe(ctx context.Context, subscribe model.Subscribe) (model.Subscribe, error) {
	subscribe, err := cS.subscribeRepository.CreateSubscribe(ctx, subscribe)
	if err != nil {
		return model.Subscribe{}, err
	}
	return subscribe, nil
}

func (cS *ClientService) Unsubscribe(ctx context.Context, token string) (int, error) {
	return 0, nil
}
