package src

import (
	"github.com/jinzhu/gorm"
	"time"
)

// Db database
var Db *gorm.DB

// AdmPass error
var AdmPass string

// Err error
var Err error

// Version of api
const Version = "0.2.1"

// DbLogMode log mode for database
const DbLogMode = true

// Port application use this port to get requests
const Port = "55555"

// UserSignIn struct to get user token
type UserSignIn struct {
	UserName string
	Password string
}

// UserRefreshToken struct to get user token
type UserRefreshToken struct {
	UserName     string
	RefreshToken string
}

// TokenResponse response to front end with the token and expiry time
type TokenResponse struct {
	UserID             uint
	Token              string
	TokenExpire        time.Time
	RefreshToken       string
	RefreshTokenExpire time.Time
}
