package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Shyp/go-dberror"
	"github.com/flatropw/gopher-twitter/internal/app/models"
)


func Connect(connString string) (*sql.DB, error) {
	return sql.Open("postgres", connString)
}

type UsersRepositoryPostgres struct {
	db *sql.DB
}

func NewUsersRepositoryPostgres(db *sql.DB) *UsersRepositoryPostgres {
	return &UsersRepositoryPostgres{db: db}
}

func (r *UsersRepositoryPostgres) Save(user models.User) (models.User, error) {
	err := r.db.QueryRow(InsertQuery, user.Login, user.Email, user.Password, user.Token).Scan(&user.Id)
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return models.User{}, fmt.Errorf(e.Error())
	default:
		return user, nil
	}
}

func (r *UsersRepositoryPostgres) ListAll() (users []models.User, err error) {
	rows, err := r.db.Query(ListAllQuery)
	defer func() {
		_ = rows.Close()
	}()
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return users, fmt.Errorf(e.Error())
	default:
		for rows.Next() {
			var user models.User
			err = rows.Scan(&user.Id, &user.Login, &user.Email, &user.Password, &user.Token)
			if err != nil {
				return
			}
			users = append(users, user)
		}
	}

	return
}

func (r *UsersRepositoryPostgres) GetByID(id uint) (user models.User, err error) {
	row := r.db.QueryRow(GetByIdQuery, id)
	err = row.Scan(&user.Id, &user.Login, &user.Email, &user.Password, &user.Token)
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return models.User{}, fmt.Errorf(e.Error())
	default:
		return
	}
}


func (r *UsersRepositoryPostgres) GetByEmail(email string) (user models.User, err error) {
	row := r.db.QueryRow(GetByEmailQuery, email)
	err = row.Scan(&user.Id, &user.Login, &user.Email, &user.Password, &user.Token)
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return models.User{}, fmt.Errorf(e.Error())
	default:
		return
	}
}

func (r *UsersRepositoryPostgres) GetByLogin(login string) (user models.User, err error) {
	row := r.db.QueryRow(GetByLoginQuery, login)
	err = row.Scan(&user.Id, &user.Login, &user.Email, &user.Password, &user.Token)
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return models.User{}, fmt.Errorf(e.Error())
	default:
		return
	}
}


func (r *UsersRepositoryPostgres) Delete(id uint) error {
	_, err := r.db.Exec(DeleteQuery, id)
	return dberror.GetError(err)
}