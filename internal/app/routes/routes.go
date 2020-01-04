package routes

import (
	"github.com/flatropw/gopher-twitter/internal/app/controllers"
	"github.com/flatropw/gopher-twitter/internal/app/middlewares"
	"github.com/gorilla/mux"
)

var Router = mux.NewRouter()

func Init() {
	Router.Use(middlewares.JwtAuthentication)
	Router.HandleFunc("/", controllers.Index).Methods("GET")
	Router.HandleFunc("/api/v1/user/register", controllers.RegisterUser).Methods("POST")
	Router.HandleFunc("/api/v1/user/login", controllers.Authenticate).Methods("POST")
}

