package main

import (
	"math/rand"
	"net/http"
	"strconv"
)

type HelloHandler struct{}

func NewHelloHandler(router *http.ServeMux) {
	handler := &HelloHandler{}
	router.HandleFunc("/random", handler.RollDice())
}

func (handler *HelloHandler) RollDice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		result := rand.Intn(6) + 1

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(strconv.Itoa(result)))
	}

}
