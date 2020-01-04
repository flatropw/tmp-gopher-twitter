package models

import (
	"fmt"
	"github.com/Shyp/go-dberror"
	"github.com/flatropw/gopher-twitter/internal/app/db"
	u "github.com/flatropw/gopher-twitter/internal/app/utils"
)

type Tweet struct {
	Id uint `json:"id"`
	Message string `json:"message"`
	UserId string `json:"user_id"`
	CreatedAt int64 `json:"created_at"`
}

const (
	MinTweetLength = 1
	MaxTweetLength = 280
)

func (tweet *Tweet) Validate() (map[string] interface{}, bool) {
	if len(tweet.Message) < MinTweetLength {
		return u.Message(false, fmt.Sprintf("Tweet length must be longer then %d characters", MinTweetLength)), false
	}

	if len(tweet.Message) > MaxTweetLength {
		return u.Message(false, fmt.Sprintf("Tweet length must be shorter then %d characters", MaxTweetLength)), false
	}

	return u.Message(true, "*ok*"), true
}

func (tweet *Tweet) Create() (*Tweet, error) {
	//tweet.CreatedAt = time.Now().Unix()
	err := db.Instance.Db.QueryRow(db.TweetInsertQuery, tweet.Message, tweet.UserId, tweet.CreatedAt).Scan(&tweet.Id)
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return &Tweet{}, fmt.Errorf(e.Error())
	default:
		return tweet, nil
	}
}