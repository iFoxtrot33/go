package verify

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"net/smtp"
	"validation/config"

	"github.com/jordan-wright/email"
)

type Handler struct {
	cfg           *config.Config
	verifications sync.Map
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

func (h *Handler) SendVerificationRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hash := h.generateVerificationHash(req.Email)
	verificationLink := fmt.Sprintf("http://%s/verify/%s",
		h.cfg.ServerAddress, hash)

	h.verifications.Store(hash, req.Email)

	if err := h.sendVerificationEmail(req.Email, verificationLink); err != nil {
		http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(VerificationResponse{
		Success: true,
		Message: "Verification email sent",
	})
}

func (h *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hash := strings.TrimPrefix(r.URL.Path, "/verify/")

	email, ok := h.verifications.Load(hash)
	if !ok {
		http.Error(w, "Invalid verification link", http.StatusNotFound)
		return
	}

	h.verifications.Delete(hash)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(VerificationResponse{
		Success: true,
		Message: fmt.Sprintf("Email %v verified successfully", email),
	})
}

func (h *Handler) sendVerificationEmail(toEmail, verificationLink string) error {
	e := email.NewEmail()
	e.From = h.cfg.SMTP.From
	e.To = []string{toEmail}
	e.Subject = "Email Verification"
	e.HTML = []byte(fmt.Sprintf(`
        <h1>Email Verification</h1>
        <p>Please click the link below to verify your email:</p>
        <a href="%s">Verify Email</a>
    `, verificationLink))

	auth := smtp.PlainAuth("",
		h.cfg.SMTP.From,
		h.cfg.SMTP.Password,
		h.cfg.SMTP.Host,
	)

	return e.Send(
		fmt.Sprintf("%s:%s", h.cfg.SMTP.Host, h.cfg.SMTP.Port),
		auth,
	)
}

func (h *Handler) generateVerificationHash(email string) string {
	hasher := sha256.New()
	hasher.Write([]byte(email + time.Now().String()))
	return hex.EncodeToString(hasher.Sum(nil))
}
