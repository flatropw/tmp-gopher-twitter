package routes

import (
	"github.com/flatropw/gopher-twitter/internal/app/middlewares"
	"github.com/gorilla/mux"
)

func Init() {
	router := mux.NewRouter()
	router.Use(middlewares.JwtAuthentication)
}
