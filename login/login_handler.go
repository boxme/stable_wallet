package login

import (
	"net/http"
	"stable_wallet/main/internal/app"
	"stable_wallet/main/internal/validator"
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
	// Initialize a new Validator instance.
	v := validator.New()

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

		email := r.PostForm.Get("email")
		v.Check(strings.TrimSpace(email) != "", "email", "Email is not provided")
		v.Check(utf8.RuneCountInString(email) <= 100, "email", "Email is too long")
		v.Check(validator.Matches(email, validator.EmailRX), "email", "Email is not valid")

		mobileNumber := r.PostForm.Get("mobile_number")
		v.Check(strings.TrimSpace(mobileNumber) != "", "mobile_number", "Mobile number is not provided")

		password := r.PostForm.Get("password")
		v.Check(strings.TrimSpace(password) != "", "password", "Password is not provided")
		v.Check(utf8.RuneCountInString(password) <= 100, "password", "Password is too long")

		if !v.Valid() {
			lh.app.FailedValidationResponse(w, r, v.Errors)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		lh.app.InfoLog.Printf("handling user login successful")
	}
}
