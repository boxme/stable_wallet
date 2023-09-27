package login

import (
	"context"
	"stable_wallet/main/internal/app"
	"stable_wallet/main/internal/data"
)

type LoginService interface {
	Login(ctx context.Context, email string, countryCode int, mobileNumber string, passwordPlaintext string) (*data.User, error)
	Signup(ctx context.Context, email string, countryCode int, mobileNumber string, passwordPlaintext string) (*data.User, error)
}

type loginService struct {
	app *app.App
}

func CreateLoginService(app *app.App) LoginService {
	return &loginService{
		app: app,
	}
}

func (ls *loginService) Signup(
	ctx context.Context,
	email string,
	countryCode int,
	mobileNumber string,
	passwordPlaintext string) (*data.User, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		user, err := data.Signup(*ls.app, ctx, email, countryCode, mobileNumber, passwordPlaintext)
		if err != nil {
			return nil, err
		}

		return user, nil
	}
}

func (ls *loginService) Login(
	ctx context.Context,
	email string,
	countryCode int,
	mobileNumber string,
	passwordPlaintext string) (*data.User, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		user, err := data.Login(*ls.app, ctx, email, countryCode, mobileNumber, passwordPlaintext)
		if err != nil {
			return nil, err
		}

		return user, nil
	}
}
