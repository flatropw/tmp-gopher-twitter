package db

//users queries
const (
	InsertQuery     = "INSERT INTO users (login, email, password, token) VALUES ($1, $2, $3, $4) RETURNING id;"
	ListAllQuery    = "SELECT id, login, email, password, token FROM users;"
	GetByIdQuery    = "SELECT id, login, email, password, token FROM users WHERE id = $1;"
	GetByEmailQuery = "SELECT id, login, email, password, token FROM users WHERE email = $1;"
	GetByLoginQuery = "SELECT id, login, email, password, token FROM users WHERE login = $1;"
	DeleteQuery     = "DELETE FROM users WHERE id = $1;"
)

//tweet queries

const (
	TweetInsertQuery = "INSERT INTO tweets (message, user_id, created_at) VALUES ($1, $2, $3) RETURNING id;"
)
