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
		UserId: r.Context().Value("user") . (uint),
		CreatedAt: 	time.Now().Unix(),
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
