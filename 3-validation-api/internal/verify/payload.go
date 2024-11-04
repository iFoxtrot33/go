package verify

import "time"

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
}

type SendRequest struct {
	Email string `json:"email"validate:"required,email"`
}

type SendResponse struct {
	Response string `json:"response"`
}

type RecoveryData struct {
	Hash      string    `json:"hash"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type RecoveryRequest struct {
	Hash string `json:"hash"`
}

type RecoveryResponse struct {
	Response string `json:"response"`
}
