package product

import (
	"net/http"
	"order-api/pkg/req"
	"order-api/pkg/res"
	"strconv"

	"gorm.io/gorm"
)

type ProductHandlerDeps struct {
	ProductRepository *ProductRepository
}

type ProductHandler struct {
	ProductRepository *ProductRepository
}

func NewOrderHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	handler := &ProductHandler{
		ProductRepository: deps.ProductRepository,
	}

	router.HandleFunc("POST /product", handler.Create())
	router.HandleFunc("PATCH /product/{id}", handler.Update())
	router.HandleFunc("DELETE /product/{id}", handler.Delete())
	router.HandleFunc("GET /product/{id}", handler.GetById())
	router.HandleFunc("GET /product/all", handler.GetAll())
}

func (handler *ProductHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[ProductCreateRequest](&w, r)

		if err != nil {
			return
		}
		product := NewProduct(body)
		createdProduct, err := handler.ProductRepository.Create(product)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, createdProduct, 201)
	}
}

func (handler *ProductHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[ProductCreateRequest](&w, r)

		if err != nil {
			return
		}

		idString := r.PathValue("id")

		if idString == "" {
			http.Error(w, "no id", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseUint(idString, 10, 32)

		if err != nil {
			http.Error(w, "invalid id format", http.StatusBadRequest)
			return
		}

		product, err := handler.ProductRepository.Update(&Product{
			Model:       gorm.Model{ID: uint(id)},
			Name:        body.Name,
			Description: body.Description,
			Price:       body.Price,
			Images:      body.Images,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, product, 201)
	}
}

func (handler *ProductHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = handler.ProductRepository.GetById(uint(id))

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		err = handler.ProductRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, nil, 200)
	}
}

func (handler *ProductHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")

		if idString == "" {
			http.Error(w, "no id", http.StatusBadRequest)
			return
		}

		idUint, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, "invalid id format", http.StatusBadRequest)
			return
		}

		product, err := handler.ProductRepository.GetById(uint(idUint))

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		res.Json(w, product, http.StatusOK)
	}
}

func (handler *ProductHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := handler.ProductRepository.GetAll()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, products, http.StatusOK)

	}
}
