package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "strconv"

	"github.com/gorilla/mux"
)

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
		Db.Find(&users)

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

// UserOptionGetPatchDelete crud for user object
func UserOptionGetPatchDelete(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	fmt.Println(params)
	fmt.Printf("Req: %s %s\n", r.Host, r.URL.Path)

	switch r.Method {
	case "OPTIONS":
		log.Println("OPTIONS /user/")

	case "PATCH":
		log.Println("PATCH /user/" + params["id"])

	case "DELETE":
		log.Println("DELETE /user/" + params["id"])

	case "GET":

		log.Println("GET /user/" + params["id"])
		// id, _ := strconv.ParseUint(params["id"], 10, 64)
		var user User
		Db.Where("id = ?", params["id"]).Last(&user)
		if user.Name != "" {
			addedrecordString, _ := json.Marshal(user)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, string(addedrecordString))
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "{\"message\":\"User Not Found\"}")

		}

	default:
		fmt.Fprintf(w, "Sorry, only PATCH,DELETE,GET methods are supported.")
	}
}

// UserPostOptions post new user
func UserPostOptions(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "OPTIONS":
		log.Println("OPTIONS /user")
	case "POST":
		log.Println("POST /user")
	default:
		fmt.Fprintf(w, "Sorry, only OPTIONS,POST methods are supported.")
	}
}
