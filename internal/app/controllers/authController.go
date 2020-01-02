package controllers

import (
	"encoding/json"
	"github.com/flatropw/gopher-twitter/internal/app/models"
	u "github.com/flatropw/gopher-twitter/internal/app/utils"
	"io/ioutil"
	"net/http"
)

var RegisterUser = func(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}

	response := user.Create()
	u.Response(w, response)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}

	response := models.Login(user.Email, user.Password)
	u.Response(w, response)
}