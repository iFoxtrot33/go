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

func App() http.Handler {

	conf := configs.LoadConfig()

	database := db.NewDb(conf)

	router := http.NewServeMux()

	//repositories
	productRepository := product.NewProductRepository(database)
	userRepository := user.NewUserRepository(database)
	orderRepository := order.NewOrderRepository(database)

	//services
	authService := auth.NewAuthService(userRepository)

	//handlers
	product.NewProductHandler(router, product.ProductHandlerDeps{
		ProductRepository: productRepository,
	})

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	order.NewOrderHandler(router, order.OrderHandlerDeps{
		OrderRepository: orderRepository,
		Config:          conf,
	})

	//middlewares
	stack := middleware.Chain(
		middleware.Logging,
		middleware.TokenMiddleware(conf.Auth.Secret),
	)

	return stack(router)

}

func main() {

	app := App()

	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	logrus.Info("Server starting at http://localhost:8081")
	if err := server.ListenAndServe(); err != nil {
		logrus.Fatalf("Server failed to start: %v", err)
	}
}
