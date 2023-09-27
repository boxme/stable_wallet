package data

import (
	"context"
	"errors"
	"net/http"
	"net/mail"
	"stable_wallet/main/internal/app"
	"stable_wallet/main/internal/validator"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id          uint64    `json:"id"`
	Email       string    `json:"email"`
	Token       string    `json:"-"` // Use the - directive
	CreatedAt   time.Time `json:"created_at"`
	Activated   bool      `json:"activated"`
	CountryCode int       `json:"country_code"`
	Password    password  `json:"-"`
	PhoneNumber string    `json:"phone_number"`
}

type password struct {
	plainttext *string
	hash       []byte
}

type JwtClaims struct {
	MobileNumber string
	jwt.RegisteredClaims
}

func Signup(
	app app.App,
	ctx context.Context,
	email string,
	countryCode int,
	mobileNumber string,
	passwordPlaintext string) (*User, error) {
	password, err := Create(passwordPlaintext)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO users (phone_number, country_code, password_hash, email, activated) VALUES (@phoneNumber, @countryCode, @passwordHash, @email, true) RETURNING id;`
	args := pgx.NamedArgs{
		"phoneNumber":  mobileNumber,
		"countryCode":  countryCode,
		"passwordHash": password.hash,
		"email":        email}

	var userId uint64
	err = app.Db.QueryRow(ctx, query, args).Scan(&userId)
	if err != nil {
		return nil, err
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claim := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}
	jwtClaim := &JwtClaims{MobileNumber: mobileNumber, RegisteredClaims: claim}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)
	tokenString, err := token.SignedString(app.JwtSecretKey)
	if err != nil {
		return nil, err
	}

	return &User{
		Id:          userId,
		Email:       email,
		Token:       tokenString,
		Activated:   true,
		CountryCode: countryCode,
		PhoneNumber: mobileNumber}, nil
}

func Login(
	app app.App,
	ctx context.Context,
	email string,
	countryCode int,
	mobileNumber string,
	passwordPlaintext string) (*User, error) {

	query := `SELECT id, password_hash FROM users WHERE phone_number = @phoneNumber AND country_code = @countryCode;`
	args := pgx.NamedArgs{
		"phoneNumber": mobileNumber,
		"countryCode": countryCode}

	var userId uint64
	var hashedPassword string
	err := app.Db.QueryRow(ctx, query, args).Scan(&userId, &hashedPassword)
	if err != nil {
		return nil, err
	}

	passwordMatched, err := Matches(hashedPassword, passwordPlaintext)
	if err != nil {
		return nil, err
	}

	if !passwordMatched {
		return nil, errors.New("User is not found")
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claim := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}
	jwtClaim := &JwtClaims{MobileNumber: mobileNumber, RegisteredClaims: claim}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)
	tokenString, err := token.SignedString(app.JwtSecretKey)
	if err != nil {
		return nil, err
	}

	return &User{
		Id:          userId,
		Email:       email,
		Token:       tokenString,
		Activated:   true,
		CountryCode: countryCode,
		PhoneNumber: mobileNumber}, nil
}

func Create(plaintextPassword string) (password, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	var newPassword password
	if err != nil {
		return newPassword, err
	}

	newPassword = password{plainttext: &plaintextPassword, hash: hash}
	return newPassword, nil
}

func Matches(hashedPassword string, plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil

		default:
			return false, err
		}
	}
	return true, nil

}

func (p *password) Set(plaintextPassword string) error {
	// Cost of 12
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plainttext = &plaintextPassword
	p.hash = hash

	return nil
}

func RenewJwt(app app.App, w http.ResponseWriter, r *http.Request) (string, bool) {
	token, claims := retrieveToken(app, w, r)
	if token != nil && claims != nil {
		return "", false
	}

	// New token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return "", false
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenString, err := newToken.SignedString(app.JwtSecretKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "", false
	}

	return newTokenString, true
}

/*
 * Verify JWT token from http request header. Return true if token is valid.
 */
func ValidateJwt(app app.App, w http.ResponseWriter, r *http.Request) bool {
	token, claims := retrieveToken(app, w, r)

	return token != nil && claims != nil
}

func retrieveToken(app app.App, w http.ResponseWriter, r *http.Request) (*jwt.Token, *JwtClaims) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return nil, nil
		}

		w.WriteHeader(http.StatusBadRequest)
		return nil, nil
	}

	tokenString := cookie.Value
	claims := &JwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return app.JwtSecretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader((http.StatusUnauthorized))
			return nil, nil
		}
		w.WriteHeader((http.StatusBadRequest))
		return nil, nil
	}

	if !token.Valid {
		w.WriteHeader(http.StatusBadRequest)
		return nil, nil
	}

	return token, claims
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(strings.TrimSpace(email) != "", "email", "Email is not provided")
	v.Check(utf8.RuneCountInString(email) <= 100, "email", "Email is too long")
	_, err := mail.ParseAddress(email)
	v.Check(err == nil, "email", "Email is not valid")
}

func ValidatePlaintextPassword(v *validator.Validator, password string) {
	v.Check(strings.TrimSpace(password) != "", "password", "Password is not provided")
	v.Check(utf8.RuneCountInString(password) <= 100, "password", "Password is too long")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateContactNumber(v *validator.Validator, mobileNumber string, countryCode int) {
	v.Check(countryCode > 0, "country_code", "Country code is invalid")
	v.Check(strings.TrimSpace(mobileNumber) != "", "mobile_number", "Mobile number is not provided")
}
