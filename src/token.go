package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

// Token Described Roles
type Token struct {
	gorm.Model
	Token   string `gorm:"unique_index"`
	UserID  uint
	RoleID  uint
	Expired time.Time
}

// RefreshToken Described Roles
type RefreshToken struct {
	gorm.Model
	RefreshToken string `gorm:"unique_index"`
	UserID       uint
	RoleID       uint
	Expired      time.Time
}

// GetToken get token Sign In Procedure
func GetToken(w http.ResponseWriter, r *http.Request) {

	FillAnswerHeader(w)
	OptionsAnswer(w)

	switch r.Method {

	case "POST":

		log.Println("POST /token")
		var usi UserSignIn
		err := json.NewDecoder(r.Body).Decode(&usi)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var currentUser User
		Db.Where("user_name = ?", usi.UserName).Last(&currentUser)
		if currentUser.UserName == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "{\"message\":\"User not found\"}")
			return
		}

		if !currentUser.Enabled {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "{\"message\":\"User is not active\"}")
			return
		}

		if comparePasswords(currentUser.Hash, []byte(usi.Password)) {

			now := time.Now()

			var TokenRec Token
			TokenRec.Token = GetHash()
			TokenRec.UserID = currentUser.ID
			TokenRec.RoleID = currentUser.RoleID
			TokenRec.Expired = now.AddDate(0, 0, 7)
			Db.Create(&TokenRec)

			var RTokenRec RefreshToken
			RTokenRec.RefreshToken = GetHash()
			RTokenRec.UserID = currentUser.ID
			RTokenRec.RoleID = currentUser.RoleID
			RTokenRec.Expired = now.AddDate(0, 0, 14)
			Db.Create(&RTokenRec)

			var TR TokenResponse
			TR.UserID = currentUser.ID
			TR.Token = TokenRec.Token
			TR.TokenExpire = TokenRec.Expired
			TR.RefreshToken = RTokenRec.RefreshToken
			TR.RefreshTokenExpire = RTokenRec.Expired

			apiTokenResponse, _ := json.Marshal(TR)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, string(apiTokenResponse))

			log.Println("POST /token DONE")

		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "{\"message\":\"Wrong password\"}")
			return
		}

	default:
		fmt.Fprintf(w, "Sorry, only POST method are supported.")
	}

}

// GetRefreshToken get token by refresh token Procedure
func GetRefreshToken(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		log.Println("POST /refreshtoken")

	default:
		fmt.Fprintf(w, "Sorry, only POST method are supported.")
	}
}
