package verify

type VerificationRequest struct {
	Email string `json:"email"`
}

type VerificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
