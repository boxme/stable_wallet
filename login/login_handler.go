package login

import (
	"fmt"
	"net/http"
	app "stable_wallet/main/internal"
	"strings"
	"unicode/utf8"
)

type LoginHandler struct {
	app     *app.App
	service LoginService
}

func CreateLoginHandler(app *app.App) *LoginHandler {
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
			lh.app.MethodNotAllowed(w, r)
			return
		}

		if err := r.ParseForm(); err != nil {
			lh.app.BadRequestResponse(w, r, err)
			return
		}

		// login
		errors := make(map[string]string)

		email := r.PostForm.Get("email")
		if strings.TrimSpace(email) == "" {
			errors["email"] = "Email is not provided"
		} else if utf8.RuneCountInString(email) > 100 {
			errors["email"] = "Email is too long"
		}

		password := r.PostForm.Get("password")
		if strings.TrimSpace(password) == "" {
			errors["password"] = "Password is not provided"
		} else if utf8.RuneCountInString(password) > 100 {
			errors["password"] = "Password is too long"
		}

		if len(errors) > 0 {
			fmt.Fprint(w, errors)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		lh.app.InfoLog.Printf("handling user login successful")
	}
}
