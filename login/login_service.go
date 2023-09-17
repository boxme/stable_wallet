package login

import (
	"context"
	"stable_wallet/main/internal/app"
	"stable_wallet/main/internal/data"
)

type LoginService interface {
	Login(ctx context.Context, mobileNumber string, passwordPlaintext string) (*data.User, error)
}

type loginService struct {
	app *app.App
}

func CreateLoginService(app *app.App) LoginService {
	return &loginService{
		app: app,
	}
}

func (ls *loginService) Login(
	ctx context.Context, mobileNumber string, passwordPlaintext string) (*data.User, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		password, err := data.Create(passwordPlaintext)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}
