package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Phone     string `json:"phone" gorm:"unique"`
	Name      string
	SessionId string `gorm:"column:session_id"`
}
