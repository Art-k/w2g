package src

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Db database
var Db *gorm.DB

// AdmPass error
var AdmPass string
var DEV bool

// Err error
var Err error

// Version of api
const Version = "0.2.2"

// DbLogMode log mode for database
const DbLogMode = false

// Port application use this port to get requests
const Port = "55555"

// UserSignIn struct to get user token
type UserSignIn struct {
	Name     string
	Password string
}

// UserRefreshToken struct to get user token
type UserRefreshToken struct {
	UserName     string
	RefreshToken string
}

//Model base model for all struct
type Model struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt *time.Time
	// DeletedBy string
}

// BeforeCreate create id
func (base *Model) BeforeCreate(scope *gorm.Scope) error {
	// uuID, err := uuid.NewRandom()
	// if err != nil {
	// 	return err
	// }
	return scope.SetColumn("id", GetHash())
}

// Token Described Roles
type Token struct {
	Model
	Token   string `gorm:"unique_index"`
	UserID  string
	RoleID  string
	Expired time.Time
}

// RefreshToken Described Roles
type RefreshToken struct {
	Model
	RefreshToken string `gorm:"unique_index"`
	UserID       string
	RoleID       string
	Expired      time.Time
}

// TokenResponse response to front end with the token and expiry time
type TokenResponse struct {
	UserID             string
	Token              string
	TokenExpire        time.Time
	RefreshToken       string
	RefreshTokenExpire time.Time
}

// UserRoles the table where User Project and his role for the project is linked
type UserRoles struct {
	Model
	RoleID string
	UserID string
	TeamID string
}

// User User Data
type User struct {
	Model
	Name       string `gorm:"type:varchar(100);unique_index"`
	FullName   string
	Email      string
	RoleID     string
	Salt       string `json:"-"`
	Hash       string `json:"-"`
	SetPass    string `json:"-"`
	Enabled    bool
	Active     bool `gorm:"default:'false'"`
	PwdChanged *time.Time
}

// Role Described Roles
type Role struct {
	Model
	Name    string `gorm:"type:varchar(100);unique_index"`
	Users   []User `gorm:"foreignkey:RoleID;association_foreignkey:id"`
	GroupID string
}

// RolePermission the permission connected to Role
type RolePermission struct {
	Model
	Post     bool `gorm:"default:'false'"`
	Get      bool `gorm:"default:'false'"`
	Patch    bool `gorm:"default:'false'"`
	Delete   bool `gorm:"default:'false'"`
	EndPoint string
	RoleID   string `gorm:"foreignkey:RoleID"`
}

// UserPermission the permission connected to User
type UserPermission struct {
	Model
	Post     bool `gorm:"default:'false'"`
	Get      bool `gorm:"default:'false'"`
	Patch    bool `gorm:"default:'false'"`
	Delete   bool `gorm:"default:'false'"`
	EndPoint string
	UserID   string `gorm:"foreignkey:UserID"`
}

// Group the first level of groups
type Group struct {
	Model
	Name  string `gorm:"unique_index"`
	Roles []Role
}
