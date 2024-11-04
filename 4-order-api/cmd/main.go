package main

import (
	"fmt"
	"order-api/configs"

	"net/http"
	"order-api/pkg/db"
)

func main() {
	conf := configs.LoadConfig()

	_ = db.NewDb(conf)

	router := http.NewServeMux()

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server started at http://localhost:8081")
	server.ListenAndServe()
}
