package auth

import (
	"net/http"
	"order-api/configs"
	"order-api/pkg/jwt"
	"order-api/pkg/req"
	"order-api/pkg/res"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

type AuthHandler struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
	router.HandleFunc("POST /auth/session", handler.Session())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}

		sessionId, err := handler.AuthService.Login(body.Phone)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		data := RegisterAndLoginResponse{
			SessionId: sessionId,
		}
		res.Json(w, data, 200)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)

		if err != nil {
			return
		}

		sessionId, err := handler.AuthService.Register(body.Phone, body.Name)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		data := RegisterAndLoginResponse{
			SessionId: sessionId,
		}

		res.Json(w, data, 200)

	}
}

func (handler *AuthHandler) Session() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[AuthRequest](&w, r)

		if err != nil {
			return
		}

		existedUser, _ := handler.UserRepository.FindBySessionId(body.SessionId)

		if existedUser == nil {
			http.Error(w, WrongCredentials, http.StatusUnauthorized)
			return
		}

		token, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(existedUser.Name, existedUser.Phone)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, AuthResponse{
			Token: token,
		}, http.StatusOK)
	}
}
