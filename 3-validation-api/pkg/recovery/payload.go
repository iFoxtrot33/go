package recovery

import "time"

type RecoveryData struct {
	Hash      string    `json:"hash"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
