package main

import (
	"fmt"
	"order-api/configs"
	"order-api/internal/product"
	"order-api/pkg/middleware"

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

	stack := middleware.Chain(

		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	fmt.Println("Server started at http://localhost:8081")
	server.ListenAndServe()
}
