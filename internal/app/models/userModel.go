package models

import (
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/dgrijalva/jwt-go"
	"github.com/flatropw/gopher-twitter/internal/app/db"
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

type User struct {
	Id       uint   `json:"id"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

func (user *User) Create() map[string]interface{} {
	if response, ok := user.Validate(); !ok {
		return response
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	res, err := user.Save()
	if err != nil {
		log.Print(err)
	}

	if res.Id <= 0 {
		return u.Message(false, "Failed to create user, connection error.")
	}

	tk := &Token{Id: res.Id}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(os.Getenv("token_secret")))
	if err != nil {
		log.Print(err)
	}
	user.Token = tokenString

	user.Password = ""

	response := u.Message(true, "User has been created")
	response["user"] = user
	return response

}

func (user *User) Validate() (map[string]interface{}, bool) {
	if len(user.Login) < MinLoginLength {
		return u.Message(false, fmt.Sprintf("Login length must be longer then %d characters", MinLoginLength)), false
	}

	if len(user.Password) < MinPasswordLength {
		return u.Message(false, fmt.Sprintf("Password length must be longer then %d characters", MinPasswordLength)), false
	}

	err := checkmail.ValidateFormat(user.Email)
	if err != nil {
		return u.Message(false, "Email field empty or invalid"), false
	}

	us, err := user.GetByLogin(user.Login)
	if user.Login == us.Login {
		return u.Message(false, "User with the same login already exists"), false
	}

	us, err = user.GetByEmail(user.Email)
	if user.Email == us.Email {
		return u.Message(false, "User with the same email already exists"), false
	}

	return u.Message(true, "*ok*"), true
}

func (user *User) Authenticate(email, password string) map[string]interface{} {
	user, err := user.GetByEmail(email)
	if err != nil {
		return u.Message(false, "Unknown error")
	}

	if user == nil {
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
	resp["token"] = user.Token
	return resp
}

func (user *User) Save() (*User, error) {
	err := db.Instance.Db.QueryRow(db.InsertQuery, user.Login, user.Email, user.Password, user.Token).Scan(&user.Id)
	return user, err
}

func (user *User) GetById(id uint) (*User, error) {
	us := &User{}
	row := db.Instance.Db.QueryRow(db.GetByIdQuery, id)
	err := row.Scan(&us.Id, &us.Login, &us.Email, &us.Password)
	return us, err
}

func (user *User) GetByEmail(email string) (*User, error) {
	us := &User{}
	row := db.Instance.Db.QueryRow(db.GetByEmailQuery, email)
	err := row.Scan(&us.Id, &us.Login, &us.Email, &us.Password, &us.Token)
	return us, err
}

func (user *User) GetByLogin(login string) (*User, error) {
	us := &User{}
	row := db.Instance.Db.QueryRow(db.GetByLoginQuery, login)
	err := row.Scan(&us.Id, &us.Login, &us.Email, &us.Password, &us.Token)
	return us, err
}

func (user *User) Delete(id uint) error {
	_, err := db.Instance.Db.Exec(db.DeleteQuery, id)
	return err
}
