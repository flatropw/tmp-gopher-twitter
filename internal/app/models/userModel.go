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

	usersRepository := postgres.UsersRepositoryPostgres{}
	res, err := usersRepository.Save(*user)

	if err != nil {
		log.Print(err)
	}
	
	if res.Id <= 0 {
		return u.Message(false, "Failed to create user, connection error.")
	}

	//Создать новый токен JWT для новой зарегистрированной учётной записи
	tk := &Token{Id: res.Id}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	user.Password = "" //удалить пароль

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

	usersRepository := postgres.UsersRepositoryPostgres{}
	res, err := usersRepository.GetByLogin(user.Login)
	if (User{}) != res {
		return u.Message(false, "User with the same login already exists"), false
	}

	res, err = usersRepository.GetByEmail(user.Login)
	if (User{}) != res {
		return u.Message(false, "User with the same email already exists"), false
	}

	return u.Message(true, "*ok*"), true
}

func Login(email, password string) (map[string]interface{}) {

	user := &User{}
	err := GetDB().Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Пароль не совпадает!!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Работает! Войти в систему
	user.Password = ""

	//Создать токен JWT
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString // Сохраните токен в ответе

	resp := u.Message(true, "Logged In")
	resp["user"] = user
	return resp
}

