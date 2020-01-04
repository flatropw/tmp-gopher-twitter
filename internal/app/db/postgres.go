package db

import (
	"database/sql"
)

type dbInstance struct {
	Db *sql.DB
}

var Instance = dbInstance{}

func (*dbInstance) Connect(connString string) (err error) {
	conn, err := sql.Open("postgres", connString)
	if err == nil {
		Instance = dbInstance{conn}
	}
	return
}


