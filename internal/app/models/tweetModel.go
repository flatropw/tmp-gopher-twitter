package models

import (
	"fmt"
	u "github.com/flatropw/gopher-twitter/internal/app/utils"
)

type Tweet struct {
	Id uint `json:"id"`
	Message string `json:"message"`
	UserId string `json:"user_id"`
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