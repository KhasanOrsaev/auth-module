package models

import "gorm.io/gorm"

var groupTableName = "auth_groups"

type Group struct {
	gorm.Model
	Name string `gorm:"unique"`
	Role []*Role `gorm:"many2many:user_permissions"`
}

func (Group) TableName() string {
	if schema=="" {
		return groupTableName
	}
	return schema+"."+groupTableName
}

func SetGroupTableName(s string) {
	groupTableName = s
}

