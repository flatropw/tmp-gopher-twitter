package models

import (
	"fmt"
	"github.com/Shyp/go-dberror"
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
	Id uint `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
	Password string `json:"password"`
	Token string `json:"token";sql:"-"`
}

func (user *User) Create() map[string] interface{} {
	if resp, ok := user.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	res, err := user.Save()
	fmt.Println("res", res)
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
		return u.Message(false, fmt.Sprintf("Password length must be longer then %d characters", MinPasswordLength)), false
	}

	err := checkmail.ValidateFormat(user.Email)
	if err != nil {
		return u.Message(false, "Email field empty or invalid"), false
	}

	tmp, err := user.GetByLogin(user.Login)
	fmt.Println(tmp)
	if user.Login == tmp.Login {
		return u.Message(false, "User with the same login already exists"), false
	}

	tmp2, err := user.GetByEmail(user.Email)
	if user.Email == tmp2.Email {
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
	resp["user"] = user
	return resp
}



func (user *User) Save() (*User, error) {
	err := db.Instance.Db.QueryRow(db.InsertQuery, user.Login, user.Email, user.Password, user.Token).Scan(&user.Id)
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return &User{}, fmt.Errorf(e.Error())
	default:
		return user, nil
	}
}

func (user *User) ListAll() (users []User, err error) {
	rows, err := db.Instance.Db.Query(db.ListAllQuery)
	defer func() {
		_ = rows.Close()
	}()
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return users, fmt.Errorf(e.Error())
	default:
		for rows.Next() {
			var u User
			err = rows.Scan(&u.Id, &u.Login, &u.Email, &u.Password, &u.Token)
			if err != nil {
				return
			}
			users = append(users, u)
		}
	}

	return
}

func (user *User) GetByID(id uint) (*User, error) {
	row := db.Instance.Db.QueryRow(db.GetByIdQuery, id)
	err := row.Scan(&user.Id, &user.Login, &user.Email, &user.Password, &user.Token)
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return &User{}, fmt.Errorf(e.Error())
	default:
		return user, nil
	}
}


func (user *User) GetByEmail(email string) (*User, error) {
	us := &User{}
	row := db.Instance.Db.QueryRow(db.GetByEmailQuery, email)
	err := row.Scan(&us.Id, &us.Login, &us.Email, &us.Password, &us.Token)
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return &User{}, fmt.Errorf(e.Error())
	default:
		return us, nil
	}
}

func (user *User) GetByLogin(login string) (*User, error) {
	us := &User{}
	row := db.Instance.Db.QueryRow(db.GetByLoginQuery, login)
	err := row.Scan(&us.Id, &us.Login, &us.Email, &us.Password, &us.Token)
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return &User{}, fmt.Errorf(e.Error())
	default:
		return us, nil
	}
}


func (user *User) Delete(id uint) error {
	_, err := db.Instance.Db.Exec(db.DeleteQuery, id)
	return dberror.GetError(err)
}

