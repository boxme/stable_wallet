package login

import (
	app "stable_wallet/main/internal"
	"stable_wallet/main/internal/data"
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
	return nil, nil
}
