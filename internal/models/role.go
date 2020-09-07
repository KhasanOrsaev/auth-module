package models

import "gorm.io/gorm"

var roleTableName = "auth_roles"

type Role struct {
	gorm.Model
	Name string `gorm:"unique"`
	Scope string
}

func (Role) TableName() string {
	if schema=="" {
		return roleTableName
	}
	return schema+"."+roleTableName
}

func SetRoleTableName(s string) {
	roleTableName = s
}