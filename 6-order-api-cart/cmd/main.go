package main

import (
	"order-api/configs"
	"order-api/internal/auth"
	"order-api/internal/order"
	"order-api/internal/product"
	"order-api/internal/user"
	"order-api/pkg/middleware"

	"net/http"
	"order-api/pkg/db"

	"github.com/sirupsen/logrus"
)

func main() {
	conf := configs.LoadConfig()

	database := db.NewDb(conf)

	router := http.NewServeMux()

	productRepository := product.NewProductRepository(database)
	userRepository := user.NewUserRepository(database)
	orderRepository := order.NewOrderRepository(database)

	authService := auth.NewAuthService(userRepository)

	product.NewProductHandler(router, product.ProductHandlerDeps{
		ProductRepository: productRepository,
	})

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	stack := middleware.Chain(
		middleware.Logging,
		middleware.TokenMiddleware(conf.Auth.Secret),
	)

	order.NewOrderHandler(router, order.OrderHandlerDeps{
		OrderRepository: orderRepository,
		Config:          conf,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	logrus.Info("Server starting at http://localhost:8081")
	if err := server.ListenAndServe(); err != nil {
		logrus.Fatalf("Server failed to start: %v", err)
	}
}
