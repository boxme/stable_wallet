package login

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LoginHandler struct {
	service LoginService
}

func CreateLoginHandler(db *pgxpool.Pool) *LoginHandler {
	logingService := CreateLoginService(db)
	return &LoginHandler{
		service: logingService,
	}
}

func (lh *LoginHandler) HandleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "Use POST for login", http.StatusMethodNotAllowed)
			return
		}
		// login
		w.WriteHeader(http.StatusOK)
	}
}
