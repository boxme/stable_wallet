package login

import (
	"net/http"
	"stable_wallet/main/config"
)

type LoginHandler struct {
	app     *config.App
	service LoginService
}

func CreateLoginHandler(app *config.App) *LoginHandler {
	logingService := CreateLoginService(app)
	return &LoginHandler{
		app:     app,
		service: logingService,
	}
}

func (lh *LoginHandler) HandleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lh.app.InfoLog.Printf("handling user login...")

		if r.Method != http.MethodPost {
			lh.app.ErrorLog.Printf("using wrong restful method in user login %s", r.Method)

			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "Use POST for login", http.StatusMethodNotAllowed)
			return
		}
		// login
		w.WriteHeader(http.StatusOK)
		lh.app.InfoLog.Printf("handling user login successful")
	}
}
