package login

import "github.com/jackc/pgx/v5/pgxpool"

type LoginUser struct {
	Id    uint64
	Email string
	Token string
}

type LoginService interface {
	Login(email string, hashPassword string, idemKey string) (*LoginUser, error)
}

type loginService struct {
	Db *pgxpool.Pool
}

func CreateLoginService(db *pgxpool.Pool) LoginService {
	return &loginService{
		Db: db,
	}
}

func (ls *loginService) Login(email string, hashPassword string, idemKey string) (*LoginUser, error) {
	return nil, nil
}
