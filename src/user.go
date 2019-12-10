package src

import (
	"encoding/json"
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
	UserName  string `gorm:"type:varchar(100);unique_index"`
	FullName  string
	Email     string
	RoleID    uint
	Salt      string
	Hash      string
	Enabled   bool
	CreatedBy uint
	UpdatedBy uint
	DeletedBy uint
}

type apiUsersResponse struct {
	API    string
	Total  int
	Entity []User
}

// Users User routine
func Users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Println("GET Users")

		var users []User
		Db.Select("ID, user_name, full_name, created_at, updated_at, email, role_id, enabled, created_by, updated_at, updated_by").Find(&users)

		var Response apiUsersResponse
		Response.Entity = users
		Response.API = Version
		Response.Total = len(Response.Entity)

		addedrecordString, _ := json.Marshal(Response)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(addedrecordString))

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
	fmt.Printf("Req: %s %s\n", r.Host, r.URL.Path)

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
