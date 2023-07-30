package login

import (
	"net/http"
)

type LoginHandler struct {
}

func CreateNewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (lh *LoginHandler) HandleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			w.WriteHeader(405)
			w.Write([]byte("Use POST for login"))
			return
		}
		// login
		w.WriteHeader(http.StatusOK)
	}
}
