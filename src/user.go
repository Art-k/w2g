package src

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// UserSectionName this name will be used to set permissions for users
const UserSectionName = "UserSection"

// User User Data
type User struct {
	gorm.Model
	UserName string `gorm:"type:varchar(100);unique_index"`
	FullName string
	Email    string
	RoleID   uint
	Salt     string
	Hash     string
	Enabled  bool
}

// Users User routine
func Users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Println("GET Users")

	default:
		fmt.Fprintf(w, "Sorry, only GET method are supported.")
	}
}

// UserCRUD crud for user object
func UserCRUD(w http.ResponseWriter, r *http.Request) {

	FillAnswerHeader(w)
	OptionsAnswer(w)

	params := mux.Vars(r)
	fmt.Println(params)

	switch r.Method {
	case "OPTIONS":
		log.Println("OPTIONS /user/" + params["id"])

	case "POST":
		log.Println("POST /user/" + params["id"])

	case "PATCH":
		log.Println("PATCH /user/" + params["id"])

	case "DELETE":
		log.Println("DELETE /user/" + params["id"])

	case "GET":
		log.Println("GET /user/" + params["id"])

	default:
		fmt.Fprintf(w, "Sorry, only OPTIONS,POST,PATCH,DELETE,GET methods are supported.")
	}
}
