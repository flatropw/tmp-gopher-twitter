package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/flatropw/gopher-twitter/internal/app/models"
	u "github.com/flatropw/gopher-twitter/internal/app/utils"
	"io/ioutil"
	"net/http"
)

var Subscribe = func(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}
	defer func(){
		_ = r.Body.Close()
	}()

	parsed := struct {
		Nickname string
	}{}
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}

	us := &models.User{}
	us, err = us.GetById(r.Context().Value("user").(uint))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		u.Response(w, u.Message(false, "User not found"))
	}

	sub := &models.Subscriber{User: us}
	userToSubscribe, err := us.GetByLogin(parsed.Nickname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.Response(w, u.Message(false, "Internal server error"))
	}

	subscription, err := sub.SubscribeTo(userToSubscribe.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.Response(w, u.Message(false, "Internal server error"))
	}
	var prefix string
	if subscription.Status == false {
		prefix = "un"
	}
	response := u.Message(true, fmt.Sprintf("You have successfully %ssubscribed", prefix))
	response["subscription"] = subscription
	u.Response(w, response)
}
