package login

import (
	"context"
	"encoding/json"
	"net/http"
	"stable_wallet/main/internal/app"
	"stable_wallet/main/internal/data"
	"stable_wallet/main/internal/validator"
	"strconv"
	"strings"
	"time"
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

func (lh *LoginHandler) HandleSignup() http.HandlerFunc {
	v := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		lh.app.InfoLog.Printf("handling user signup...")

		if r.Method != http.MethodPost {
			lh.app.ErrorLog.Printf("using wrong restful method in user signup %s", r.Method)

			w.Header().Set("Allow", http.MethodPost)
			lh.app.MethodNotAllowed(w, r)
			return
		}

		if err := r.ParseForm(); err != nil {
			lh.app.BadRequestResponse(w, r, err)
			return
		}

		email := r.PostForm.Get("email")
		data.ValidateEmail(v, email)

		mobileNumber := r.PostForm.Get("mobile_number")
		countryCode, err := strconv.Atoi(r.PostForm.Get("country_code"))
		if err != nil {
			lh.app.BadRequestResponse(w, r, err)
			return
		}

		data.ValidateContactNumber(v, mobileNumber, countryCode)
		v.Check(strings.TrimSpace(mobileNumber) != "", "mobile_number", "Mobile number is not provided")

		password := r.PostForm.Get("password")
		data.ValidatePlaintextPassword(v, password)

		if !v.Valid() {
			lh.app.FailedValidationResponse(w, r, v.Errors)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		user, err := lh.service.Signup(ctx, email, countryCode, mobileNumber, password)
		if err != nil {
			lh.app.BadRequestResponse(w, r, err)
			return
		}

		jsonResponse, err := json.Marshal(user)
		if err != nil {
			lh.app.BadRequestResponse(w, r, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
		lh.app.InfoLog.Printf("handling user signup successful")
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
		data.ValidateEmail(v, email)

		mobileNumber := r.PostForm.Get("mobile_number")
		countryCode, err := strconv.Atoi(r.PostForm.Get("country_code"))
		if err != nil {
			lh.app.BadRequestResponse(w, r, err)
			return
		}

		data.ValidateContactNumber(v, mobileNumber, countryCode)
		v.Check(strings.TrimSpace(mobileNumber) != "", "mobile_number", "Mobile number is not provided")

		password := r.PostForm.Get("password")
		data.ValidatePlaintextPassword(v, password)

		if !v.Valid() {
			lh.app.FailedValidationResponse(w, r, v.Errors)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		user, err := lh.service.Login(ctx, email, countryCode, mobileNumber, password)
		if err != nil {
			lh.app.BadRequestResponse(w, r, err)
			return
		}

		jsonResponse, err := json.Marshal(user)
		if err != nil {
			lh.app.BadRequestResponse(w, r, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
		lh.app.InfoLog.Printf("handling user login successful")
	}
}
