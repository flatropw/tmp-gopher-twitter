package controllers

import (
	"encoding/json"
	"github.com/flatropw/gopher-twitter/internal/app/models"
	u "github.com/flatropw/gopher-twitter/internal/app/utils"
	"io/ioutil"
	"net/http"
	"time"
)

var CreateTweet = func(w http.ResponseWriter, r *http.Request) {
	tweet := models.Tweet{
		UserId:    r.Context().Value("user").(uint),
		CreatedAt: time.Now().Unix(),
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}

	err = json.Unmarshal(body, &tweet)
	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}

	response := tweet.Create()
	u.Response(w, response)
}

var ShowTweets = func(w http.ResponseWriter, r *http.Request) {
	us := &models.User{}
	us, err := us.GetById(r.Context().Value("user").(uint))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		u.Response(w, u.Message(false, "User not found"))
	}

	sub := &models.Subscriber{User: us}
	subscriptions, err := sub.GetSubscriptions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.Response(w, u.Message(false, "Internal server error"))
	}

	var userIds []uint
	for _, v := range subscriptions {
		if v.Status == true {
			userIds = append(userIds, v.SubscribedId)
		}
	}

	t := models.Tweet{}
	tweets, err := t.GetByUserIds(userIds, 30)
	response := u.Message(true, "Tweets from your subscription list")
	response["tweets"] = tweets
	u.Response(w, response)
}