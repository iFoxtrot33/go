package main

import (
	"fmt"
	"order-api/configs"
	"order-api/internal/product"

	"net/http"
	"order-api/pkg/db"
)

func main() {
	conf := configs.LoadConfig()

	database := db.NewDb(conf)

	router := http.NewServeMux()

	productRepository := product.NewProductRepository(database)

	product.NewOrderHandler(router, product.ProductHandlerDeps{
		ProductRepository: productRepository,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server started at http://localhost:8081")
	server.ListenAndServe()
}
