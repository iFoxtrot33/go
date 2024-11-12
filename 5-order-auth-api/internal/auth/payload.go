package auth

type LoginRequest struct {
	Phone string `json:"phone" validate:"required,phone"`
}

type RegisterRequest struct {
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone" validate:"required,phone"`
}

type RegisterAndLoginResponse struct {
	SessionId string `json:"session_id"`
}

type AuthRequest struct {
	SessionId string `json:"SessionId"  validate:"required"`
	Code      string `json:"code"  validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
