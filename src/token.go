package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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

// IsLegalUser check if user exists
func IsLegalUser(Auth string) (bool, User) {

	var Answer bool
	var currentToken Token
	var currentUser User

	token := strings.Replace(Auth, "Bearer ", "", -1)

	Db.Where("token = ?", token).Last(&currentToken)
	if currentToken.Token != "" {

		if currentToken.Expired.After(time.Now()) {

			Db.Where("id = ?", currentToken.UserID).Last(&currentUser)

			if currentUser.UserName != "" {
				Answer = true
			} else {
				Answer = false
			}
		}

	}

	return Answer, currentUser

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

			apiTokenResponse, _ := json.Marshal(APITokenResponse(currentUser))
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

// APITokenResponse
func APITokenResponse(cu User) TokenResponse {
	now := time.Now()

	var TokenRec Token
	TokenRec.Token = GetHash()
	TokenRec.UserID = cu.ID
	TokenRec.RoleID = cu.RoleID
	TokenRec.Expired = now.AddDate(0, 0, 7)
	Db.Create(&TokenRec)

	var RTokenRec RefreshToken
	RTokenRec.RefreshToken = GetHash()
	RTokenRec.UserID = cu.ID
	RTokenRec.RoleID = cu.RoleID
	RTokenRec.Expired = now.AddDate(0, 0, 14)
	Db.Create(&RTokenRec)

	var TR TokenResponse
	TR.UserID = cu.ID
	TR.Token = TokenRec.Token
	TR.TokenExpire = TokenRec.Expired
	TR.RefreshToken = RTokenRec.RefreshToken
	TR.RefreshTokenExpire = RTokenRec.Expired

	return TR
}

// GetRefreshToken get token by refresh token Procedure
func GetRefreshToken(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		log.Println("POST /refreshtoken")

		type inRefreshToken struct {
			RefreshToken string
		}

		var rt inRefreshToken
		err := json.NewDecoder(r.Body).Decode(&rt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		Authorization := r.Header.Get("Authorization")
		_, cUser := IsLegalUserByRefreshToken(Authorization)


		apiTokenResponse, _ := json.Marshal(APITokenResponse(cUser))
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(apiTokenResponse))

	default:
		fmt.Fprintf(w, "Sorry, only POST method are supported.")
	}
}

// IsLegalUserByRefreshToken check if Refresh token connected to user
func IsLegalUserByRefreshToken(Auth string) (bool, User) {

	var Answer bool
	var currentToken Token
	var currentUser User

	token := strings.Replace(Auth, "Bearer ", "", -1)

	Db.Where("token = ?", token).Last(&currentToken)
	if currentToken.Token != "" {

		if currentToken.Expired.After(time.Now()) {

			Db.Where("id = ?", currentToken.UserID).Last(&currentUser)

			if currentUser.UserName != "" {
				Answer = true
			} else {
				Answer = false
			}
		}

	}

	return Answer, currentUser

}

