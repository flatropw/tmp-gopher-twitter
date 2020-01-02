package postgres

const (
	InsertQuery = "INSERT INTO users (login, email, password, token) VALUES ($1, $2, $3, $4);"
	ListAllQuery = "SELECT id, login, email, password, token FROM users;"
	GetByIdQuery = "SELECT id, login, email, password, token FROM users WHERE id = $1;"
	GetByEmailQuery = "SELECT id, login, email, password, token FROM users WHERE email = $1;"
	GetByLoginQuery = "SELECT id, login, email, password, token FROM users WHERE login = $1;"
	DeleteQuery = "DELETE FROM users WHERE id = $1;"
)
