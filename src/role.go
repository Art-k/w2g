package src

import (
	"github.com/jinzhu/gorm"
)

// Role Described Roles
type Role struct {
	gorm.Model
	Role  string `gorm:"type:varchar(100);unique_index"`
	Users []User `gorm:"foreignkey:RoleID"`
}
