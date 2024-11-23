package order

import (
	"log"
	"net/http"
	"order-api/configs"
	"order-api/pkg/middleware"
	"order-api/pkg/req"
	"order-api/pkg/res"
	"strconv"
)

type OrderHandlerDeps struct {
	OrderRepository *OrderRepository
	Config          *configs.Config
}

type OrderHandler struct {
	OrderRepository *OrderRepository
	Config          *configs.Config
}

func NewOrderHandler(router *http.ServeMux, deps OrderHandlerDeps) {
	handler := &OrderHandler{
		OrderRepository: deps.OrderRepository,
		Config:          deps.Config,
	}

	router.Handle("POST /order",
		middleware.TokenMiddleware(deps.Config.Auth.Secret)(
			http.HandlerFunc(handler.Create()),
		),
	)

	router.Handle("GET /order/{id}",
		middleware.TokenMiddleware(deps.Config.Auth.Secret)(
			http.HandlerFunc(handler.GetById()),
		),
	)

	router.Handle("GET /my-orders",
		middleware.TokenMiddleware(deps.Config.Auth.Secret)(
			http.HandlerFunc(handler.GetMyOrders()),
		),
	)
}

func (handler *OrderHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone := r.Context().Value(middleware.PhoneContextKey).(string)

		body, err := req.HandleBody[OrderCreateRequest](w, r)
		if err != nil {
			log.Println(err)
			return
		}

		order := &Order{
			Phone:       phone,
			Description: body.Description,
		}

		createdOrder, err := handler.OrderRepository.Create(order, body.Products)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println(err)
			return
		}

		res.Json(w, createdOrder, http.StatusCreated)
	}
}

func (handler *OrderHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			http.Error(w, "invalid id format", http.StatusBadRequest)
			log.Println(err)
			return
		}

		order, err := handler.OrderRepository.GetById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			log.Println(err)
			return
		}

		res.Json(w, order, http.StatusOK)
	}
}

func (handler *OrderHandler) GetMyOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone := r.Context().Value(middleware.PhoneContextKey).(string)

		orders, err := handler.OrderRepository.GetByPhone(phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		if orders == nil {
			orders = make([]Order, 0)
		}

		res.Json(w, orders, http.StatusOK)
	}

}
