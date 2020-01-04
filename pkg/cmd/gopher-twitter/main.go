package main

import (
	"github.com/flatropw/gopher-twitter/internal/app/db"
	"github.com/flatropw/gopher-twitter/internal/app/routes"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Panic("No .env file found")
	}

	connStr, exist := os.LookupEnv("postgres_connection_string")
	if !exist {
		log.Panic("empty postgres_connection_string (check env)")
	}
	err := db.Instance.Connect(connStr)
	if err != nil {
		panic(err)
		log.Panic(err)
	}

	routes.Init()
	port := os.Getenv("port")
	if port == "" {
		port = "8080"
	}

	err = http.ListenAndServe(":"+port, routes.Router)
	if err != nil {
		log.Panic(err)
	}
}
