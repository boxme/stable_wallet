package login

import (
	"context"
	"stable_wallet/main/internal/app"
	"stable_wallet/main/internal/data"
	"time"
)

type LoginService interface {
	Login(email string, hashPassword string, idemKey string) (*data.User, error)
}

type loginService struct {
	app *app.App
}

func CreateLoginService(app *app.App) LoginService {
	return &loginService{
		app: app,
	}
}

func (ls *loginService) Login(email string, hashPassword string, idemKey string) (*data.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if ctx != nil {

	}
	return nil, nil
}
