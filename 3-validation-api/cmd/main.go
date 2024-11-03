package main

import (
	"fmt"
	"log"
	"net/http"

	"validation/config"
	"validation/internal/verify"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	router := http.NewServeMux()
	verifyHandler := verify.NewHandler(cfg)

	router.HandleFunc("/send", verifyHandler.SendVerificationRequest)
	router.HandleFunc("/verify/", verifyHandler.VerifyEmail)

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	fmt.Printf("Server started at http://%s\n", cfg.ServerAddress)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed:", err)
	}
}
