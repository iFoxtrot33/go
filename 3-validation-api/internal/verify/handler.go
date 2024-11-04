package verify

import (
	"fmt"
	"net/http"
	"validation/config"
	"validation/pkg/recovery"
	"validation/pkg/req"
	"validation/pkg/res"

	"github.com/google/uuid"
)

type VerifyHandler struct {
	recoveryService *recovery.Service
}

func NewVerifyHandler(router *http.ServeMux, cfg *config.Config) {
	handler := &VerifyHandler{
		recoveryService: recovery.NewService(cfg.Recovery),
	}
	router.HandleFunc("POST /send", handler.Send())
	router.HandleFunc("GET /verify/{hash}", handler.Recovery())
}

func (h *VerifyHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[SendRequest](&w, r)
		if err != nil {
			return
		}

		hash := uuid.New().String()

		if err := h.recoveryService.SaveToFile(hash, body.Email); err != nil {
			res.Json(w, map[string]string{"error": "Failed to save data"}, http.StatusInternalServerError)
			return
		}

		if err := h.recoveryService.SendEmail(body.Email, hash); err != nil {
			res.Json(w, map[string]string{"error": "Failed to send email"}, http.StatusInternalServerError)
			return
		}

		response := SendResponse{
			Response: "Email sent successfully",
		}
		res.Json(w, response, http.StatusOK)
	}
}

func (h *VerifyHandler) Recovery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")

		recoveryData, err := h.recoveryService.GetRecoveryData(hash)
		if err != nil {
			res.Json(w, RecoveryResponse{Response: "Invalid recovery link"}, http.StatusBadRequest)
			return
		}

		if err := h.recoveryService.RemoveRecoveryData(hash); err != nil {
			res.Json(w, RecoveryResponse{Response: "Error processing recovery"}, http.StatusInternalServerError)
			return
		}

		response := RecoveryResponse{
			Response: fmt.Sprintf("Recovery successful for email: %s", recoveryData.Email),
		}
		res.Json(w, response, http.StatusOK)
	}
}
