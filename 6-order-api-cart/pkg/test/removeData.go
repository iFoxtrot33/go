package test

import (
	"order-api/internal/user"

	"gorm.io/gorm"
)

func RemoveData(db *gorm.DB) {
	db.Unscoped().Where("phone = ?", "+79993333333").Delete(&user.User{})
}
