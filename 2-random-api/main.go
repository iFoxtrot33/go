package main

import (
	"fmt"
	"net/http"
)

func main() {

	router := http.NewServeMux()
	NewHelloHandler(router)
	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server started at http://localhost:8081")
	server.ListenAndServe()
}
