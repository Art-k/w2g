package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "strconv"

	"github.com/gorilla/mux"
)

type apiRoleResponse struct {
	API    string
	Total  int
	Entity []Role
}

// Roles list of all roles
func Roles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		log.Println("OPTIONS Users")

	case "GET":
		log.Println("GET Users")

		// var role Role
		var roles []Role
		// var users []User

		// Db.Find(&roles)
		Db.Preload("Users").Find(&roles)

		var Response apiRoleResponse
		Response.Entity = roles
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

// RoleOptionGetPatchDelete crud for user object
func RoleOptionGetPatchDelete(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	fmt.Println(params)
	fmt.Printf("Req: %s %s\n", r.Host, r.URL.Path)

	switch r.Method {
	case "OPTIONS":
		log.Println("OPTIONS /role/")

	case "PATCH":
		log.Println("PATCH /role/" + params["id"])

	case "DELETE":
		log.Println("DELETE /role/" + params["id"])

	case "GET":

		log.Println("GET /role/" + params["id"])
		// id, _ := strconv.ParseUint(params["id"], 10, 64)
		var record Role
		Db.Where("id = ?", params["id"]).Last(&record)
		if record.Name != "" {
			addedrecordString, _ := json.Marshal(record)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, string(addedrecordString))
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "{\"message\":\"Role Not Found\"}")

		}

	default:
		fmt.Fprintf(w, "Sorry, only PATCH,DELETE,GET methods are supported.")
	}
}

// RolePostOptions post new user
func RolePostOptions(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "OPTIONS":
		log.Println("OPTIONS /role")
	case "POST":
		log.Println("POST /role")
	default:
		fmt.Fprintf(w, "Sorry, only OPTIONS,POST methods are supported.")
	}
}
