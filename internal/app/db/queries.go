package db

//users queries
const (
	InsertQuery     = "INSERT INTO users (login, email, password, token) VALUES ($1, $2, $3, $4) RETURNING id;"
	GetByIdQuery    = "SELECT id, login, email, password FROM users WHERE id = $1;"
	GetByEmailQuery = "SELECT id, login, email, password, token FROM users WHERE email = $1;"
	GetByLoginQuery = "SELECT id, login, email, password, token FROM users WHERE login = $1;"
	DeleteQuery     = "DELETE FROM users WHERE id = $1;"
)

//tweet queries

const (
	TweetInsertQuery       = "INSERT INTO tweets (message, user_id, created_at) VALUES ($1, $2, $3) RETURNING id;"
	TweetGetByUserIdsQuery = "SELECT id, message, user_id, created_at FROM tweets WHERE user_id = ANY($1) ORDER BY ID DESC LIMIT $2"
)

//subscriber queries

const (
	SubInsertQuery               = "INSERT INTO subscriptions (subscriber_id, subscribed_id, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, subscriber_id, subscribed_id, status"
	SubUpdateStatusQuery         = "UPDATE subscriptions SET status = $1, updated_at = $2 WHERE id = $3 RETURNING status"
	SubHasAlreadySubscribedQuery = "SELECT id, subscriber_id, subscribed_id, status FROM subscriptions WHERE subscriber_id = $1 AND subscribed_id = $2 LIMIT 1"
	SubGetActiveQuery            = "SELECT id, subscriber_id, subscribed_id, status FROM subscriptions WHERE subscriber_id = $1"
)
