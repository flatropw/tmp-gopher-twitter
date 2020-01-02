package models

import (
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/dgrijalva/jwt-go"
	"github.com/flatropw/gopher-twitter/internal/app/repositories/postgres"
	u "github.com/flatropw/gopher-twitter/internal/app/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

type Token struct {
	Id uint
	jwt.StandardClaims
}

const (
	MinLoginLength    = 3
	MinPasswordLength = 5
)

var repo = postgres.UsersRepositoryPostgres{}

type User struct {
	Id uint `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
	Password string `json:"password"`
	Token string `json:"token";sql:"-"`
}

func (user *User) Create() (map[string] interface{}) {
	if resp, ok := user.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	res, err := repo.Save(*user)

	if err != nil {
		log.Print(err)
	}
	
	if res.Id <= 0 {
		return u.Message(false, "Failed to create user, connection error.")
	}

	tk := &Token{Id: res.Id}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	user.Password = ""

	response := u.Message(true, "User has been created")
	response["user"] = user
	return response

}

func (user *User) Validate() (map[string] interface{}, bool) {
	if len(user.Login) < MinLoginLength {
		return u.Message(false, fmt.Sprintf("Login length must be longer then %d characters", MinLoginLength)), false
	}

	if len(user.Password) < MinPasswordLength {
		return u.Message(false, fmt.Sprintf("Password length must be longer then %d characters", MinLoginLength)), false
	}

	err := checkmail.ValidateFormat(user.Email)
	if err != nil {
		return u.Message(false, "Email field empty or invalid"), false
	}

	res, err := repo.GetByLogin(user.Login)
	if (User{}) != res {
		return u.Message(false, "User with the same login already exists"), false
	}

	res, err = repo.GetByEmail(user.Email)
	if (User{}) != res {
		return u.Message(false, "User with the same email already exists"), false
	}

	return u.Message(true, "*ok*"), true
}

func Login(email, password string) map[string]interface{} {
	user, err := repo.GetByEmail(email)
	if err != nil {
		return u.Message(false, "Unknown error")
	}

	if (User{}) == user {
		return u.Message(false, "User with this email does not exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid login credentials. Please try again")
	}

	user.Password = ""

	tk := &Token{Id: user.Id}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	user.Token, _ = token.SignedString([]byte(os.Getenv("token_password")))

	resp := u.Message(true, "Logged In")
	resp["user"] = user
	return resp
}

