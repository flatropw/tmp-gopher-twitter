package controllers

import (
	u "github.com/flatropw/gopher-twitter/internal/app/utils"
	"net/http"
)

var Index = func(w http.ResponseWriter, r *http.Request) {
	u.Response(w, u.Message(true, "Index"))
}
