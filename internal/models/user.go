package models

import (
	"github.com/google/uuid"
	"time"
)

var userTableName = "auth_users"
var schema = ""

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();"` // primary key by default
	Name string `gorm:"not null;unique;index:user__name__password_hash__idx"`
	PasswordHash string `gorm:"not null;index:user__name__password_hash__idx"`
	Role Role `gorm:"foreignKey:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Active bool
}

func (User) TableName() string {
	if schema=="" {
		return userTableName
	}
	return schema+"."+userTableName
}

func SetSchema(s *string) {
	schema = *s
}

func SetUserTableName(s string) {
	userTableName = s
}