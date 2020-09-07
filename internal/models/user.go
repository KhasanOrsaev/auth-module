package models

import (
	"gorm.io/gorm"
)

var userTableName = "auth_users"
var schema = ""

type User struct {
	gorm.Model
	Name string `gorm:"not null;unique;index:user__name__password_hash__idx"`
	PasswordHash string `gorm:"not null;index:user__name__password_hash__idx"`
	Role []*Role `gorm:"many2many:user_permissions"`
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