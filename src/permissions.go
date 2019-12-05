package src

import (
	"github.com/jinzhu/gorm"
)

// RolePermission the permission connected to Role
type RolePermission struct {
	gorm.Model
	Create bool
	Read   bool
	Update bool
	Delete bool
	Role   Role
}

// UserPermission the permission connected to User
type UserPermission struct {
	gorm.Model
	Create   bool
	Read     bool
	Update   bool
	Delete   bool
	UserName User
}
