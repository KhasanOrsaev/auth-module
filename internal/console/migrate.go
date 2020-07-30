package console

import (
	"auth-module/internal/models"
	"github.com/go-errors/errors"
	"github.com/jinzhu/gorm"
)

func Migrate(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is null")
	}
	db.AutoMigrate(&models.User{}, models.Role{})

	return nil
}
