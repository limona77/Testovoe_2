package controller

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	custom_errors "testovoe_2/internal/custom-errors"
	"testovoe_2/internal/model"
	"testovoe_2/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

type clientRoutes struct {
	clientService service.IClient
}
type UserResponse struct {
	User model.User `json:"user"`
}
type SubscribeResponse struct {
	Subscribe model.Subscribe `json:"subscribe"`
}

// @Summary AuthMe
// @Security JWT
// @Tags client
// @Description check auth
// @Accept json
// @Produce json
// @ID get-user
// @Success 200 {object} clientResponse "ok"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /me [get]
func (cR *clientRoutes) authMe(ctx *fiber.Ctx) error {
	path := "internal.controller.auth.getUser"
	accessToken := ctx.Get("Authorization")
	fmt.Println("accessToken", accessToken)
	if len([]rune(accessToken)) == 0 {
		return fmt.Errorf(path, "no token provided")
	}
	t := strings.Split(accessToken, " ")[1]
	tokenClaims, err := cR.clientService.VerifyToken(t)
	if err != nil {
		slog.Errorf(fmt.Errorf(path+".VerifyToken, error: {%w}", err).Error())
		return wrapHttpError(ctx, fiber.StatusUnauthorized, custom_errors.ErrUserUnauthorized.Error())
	}
	params := service.AuthParams{Email: tokenClaims.Email}
	user, err := cR.clientService.GetUserByEmail(ctx.Context(), params)
	if err != nil {
		if errors.Is(err, custom_errors.ErrUserNotFound) {
			slog.Errorf(fmt.Errorf(path+".GetUserByEmail, error: {%w}", err).Error())
			return wrapHttpError(ctx, fiber.StatusBadRequest, custom_errors.ErrUserNotFound.Error())
		}
		if errors.Is(err, custom_errors.ErrWrongCredetianls) {
			slog.Errorf(fmt.Errorf(path+".GetUserByEmail, error: {%w}", err).Error())
			return wrapHttpError(ctx, fiber.StatusBadRequest, custom_errors.ErrWrongCredetianls.Error())
		}
		slog.Errorf(fmt.Errorf(path+".GetUserByEmail, error: {%w}", err).Error())
		return wrapHttpError(ctx, fiber.StatusInternalServerError, "internal server error")
	}
	resp := UserResponse{User: user}
	err = httpResponse(ctx, fiber.StatusOK, resp)
	if err != nil {
		return wrapHttpError(ctx, fiber.StatusInternalServerError, "internal server error")
	}
	return nil
}

func (cR *clientRoutes) subscribe(ctx *fiber.Ctx) error {
	path := "internal.controller.subscribe.subscribe"
	accessToken := ctx.Get("Authorization")
	if len([]rune(accessToken)) == 0 {
		return fmt.Errorf(path, "no token provided")
	}
	t := strings.Split(accessToken, " ")[1]
	tokenClaims, err := cR.clientService.VerifyToken(t)
	if err != nil {
		slog.Errorf(fmt.Errorf(path+".VerifyToken, error: {%w}", err).Error())
		return wrapHttpError(ctx, fiber.StatusUnauthorized, custom_errors.ErrUserUnauthorized.Error())
	}
	p := ctx.Query("id")
	id, err := strconv.Atoi(p)
	if err != nil {
		slog.Errorf(fmt.Errorf(path+".VerifyToken, error: {%w}", err).Error())
		return wrapHttpError(ctx, fiber.StatusUnauthorized, err.Error())
	}
	params := model.Subscribe{UserID: tokenClaims.UserID, SubscribedToId: id}
	subscribe, err := cR.clientService.Subscribe(ctx.Context(), params)
	if err != nil {
		slog.Errorf(fmt.Errorf(path+".Subscribe, error: {%w}", err).Error())
		return err
	}
	resp := SubscribeResponse{Subscribe: subscribe}
	err = httpResponse(ctx, fiber.StatusOK, resp)
	if err != nil {
		return wrapHttpError(ctx, fiber.StatusInternalServerError, "internal server error")
	}
	return nil
}

func (cR *clientRoutes) unSubscribe(ctx *fiber.Ctx) error {
	path := "internal.controller.subscribe.unSubscribe"
	accessToken := ctx.Get("Authorization")
	if len([]rune(accessToken)) == 0 {
		return fmt.Errorf(path, "no token provided")
	}
	t := strings.Split(accessToken, " ")[1]
	tokenClaims, err := cR.clientService.VerifyToken(t)
	if err != nil {
		slog.Errorf(fmt.Errorf(path+".VerifyToken, error: {%w}", err).Error())
		return wrapHttpError(ctx, fiber.StatusUnauthorized, custom_errors.ErrUserUnauthorized.Error())
	}
	p := ctx.Query("id")
	id, err := strconv.Atoi(p)
	if err != nil {
		slog.Errorf(fmt.Errorf(path+".VerifyToken, error: {%w}", err).Error())
		return wrapHttpError(ctx, fiber.StatusUnauthorized, err.Error())
	}

	params := model.Subscribe{UserID: tokenClaims.UserID, SubscribedToId: id}
	err = cR.clientService.Unsubscribe(ctx.Context(), params)
	if err != nil {
		slog.Errorf(fmt.Errorf(path+".Unsubscribe, error: {%w}", err).Error())
		return err
	}

	resp := SubscribeResponse{Subscribe: params}
	err = httpResponse(ctx, fiber.StatusOK, resp)
	if err != nil {
		return wrapHttpError(ctx, fiber.StatusInternalServerError, "internal server error")
	}
	return nil
}
