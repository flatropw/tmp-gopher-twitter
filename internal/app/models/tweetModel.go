package models

import (
	"fmt"
	"github.com/Shyp/go-dberror"
	"github.com/flatropw/gopher-twitter/internal/app/db"
	u "github.com/flatropw/gopher-twitter/internal/app/utils"
	"github.com/lib/pq"
	"log"
)

type Tweet struct {
	Id        uint   `json:"id"`
	Message   string `json:"message"`
	UserId    uint   `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
}

const (
	MinTweetLength = 1
	MaxTweetLength = 280
)

func (tweet *Tweet) Validate() (map[string]interface{}, bool) {
	if len(tweet.Message) < MinTweetLength {
		return u.Message(false, fmt.Sprintf("Tweet length must be longer then %d characters", MinTweetLength)), false
	}

	if len(tweet.Message) > MaxTweetLength {
		return u.Message(false, fmt.Sprintf("Tweet length must be shorter then %d characters", MaxTweetLength)), false
	}

	return u.Message(true, "*ok*"), true
}

func (tweet *Tweet) Create() map[string]interface{} {
	if resp, ok := tweet.Validate(); !ok {
		return resp
	}

	res, err := tweet.Save()
	if err != nil {
		log.Print(err)
	}

	if res.Id <= 0 {
		return u.Message(false, "Failed to save tweet, connection error.")
	}

	response := u.Message(true, "Tweet has been saved")
	response["tweet"] = tweet
	return response
}

func (tweet *Tweet) Save() (*Tweet, error) {
	err := db.Instance.Db.QueryRow(db.TweetInsertQuery, tweet.Message, tweet.UserId, tweet.CreatedAt).Scan(&tweet.Id)
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return &Tweet{}, fmt.Errorf(e.Error())
	default:
		return tweet, nil
	}
}

func (tweet *Tweet) GetByUserId(userId uint, limit uint) (tweets []*Tweet, err error) {
	rows, err := db.Instance.Db.Query("SELECT id, message, user_id, created_at FROM tweets WHERE user_id = $1 LIMIT $2", userId, limit)
	defer rows.Close()
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return tweets, fmt.Errorf(e.Error())
	default:
	}

	for rows.Next() {
		tmp := Tweet{}
		err = rows.Scan(&tmp.Id, &tmp.Message, &tmp.UserId, &tmp.CreatedAt)
		dbErr := dberror.GetError(err)
		switch e := dbErr.(type) {
		case *dberror.Error:
			return tweets, fmt.Errorf(e.Error())
		default:
			tweets = append(tweets, &tmp)
		}
	}

	err = rows.Err()
	dbErr = dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return tweets, fmt.Errorf(e.Error())
	default:
	}
	return
}

func (tweet *Tweet) GetByUserIds(subIds []uint, limit uint) (tweets []*Tweet, err error) {
	rows, err := db.Instance.Db.Query(db.TweetGetByUserIdsQuery, pq.Array(subIds), limit)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return tweets, fmt.Errorf(e.Error())
	default:
	}

	for rows.Next() {
		tmp := Tweet{}
		err = rows.Scan(&tmp.Id, &tmp.Message, &tmp.UserId, &tmp.CreatedAt)
		dbErr := dberror.GetError(err)
		switch e := dbErr.(type) {
		case *dberror.Error:
			return tweets, fmt.Errorf(e.Error())
		default:
			tweets = append(tweets, &tmp)
		}
	}

	err = rows.Err()
	dbErr = dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return tweets, fmt.Errorf(e.Error())
	default:
	}
	return
}

