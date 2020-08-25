package migration

import (
	"auth-module/internal/models"
	"github.com/go-errors/errors"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is null")
	}
	db.AutoMigrate(&models.User{}, &models.Role{})

	return nil
}
