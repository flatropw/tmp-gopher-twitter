package controllers

import (
	"encoding/json"
	"github.com/flatropw/gopher-twitter/internal/app/models"
	u "github.com/flatropw/gopher-twitter/internal/app/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var CreateTweet = func(w http.ResponseWriter, r *http.Request) {
	tweet := models.Tweet{
		UserId:    r.Context().Value("user").(uint),
		CreatedAt: time.Now().Unix(),
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

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

	limit, err := strconv.ParseUint(r.URL.Query()["limit"][0], 10, 32)
	if err != nil {
		limit = models.TweetsLimit
	}

	tweets, err := t.GetByUserIds(userIds, limit)
	response := u.Message(true, "Tweets from your subscription list")
	response["tweets"] = tweets
	response["limit"] = limit
	u.Response(w, response)
}
