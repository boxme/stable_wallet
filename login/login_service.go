package login

type LoginUser struct {
	Id    uint64
	Email string
	Token string
}

type LoginService interface {
	Login(email string, hashPassword string, idemKey string) (*LoginUser, error)
}

type loginService struct {
}
