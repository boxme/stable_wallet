package login

import (
	"stable_wallet/main/config"
)

type LoginUser struct {
	Id    uint64
	Email string
	Token string
}

type LoginService interface {
	Login(email string, hashPassword string, idemKey string) (*LoginUser, error)
}

type loginService struct {
	app *config.App
}

func CreateLoginService(app *config.App) LoginService {
	return &loginService{
		app: app,
	}
}

func (ls *loginService) Login(email string, hashPassword string, idemKey string) (*LoginUser, error) {
	return nil, nil
}
