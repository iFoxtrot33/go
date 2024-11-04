package main

import (
	"fmt"
	"net/http"
	"validation/internal/verify"

	"validation/config"
)

func main() {

	router := http.NewServeMux()

	cfg := config.Load()

	verify.NewVerifyHandler(router, cfg)

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server started at http://localhost:8081")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
