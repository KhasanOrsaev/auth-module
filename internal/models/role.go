package models

import "github.com/lib/pq"

var roleTableName = "auth_roles"

type Role struct {
	ID int
	Name string
	Scopes pq.StringArray `gorm:"type:json"`
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