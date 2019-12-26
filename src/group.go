package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// "github.com/google/uuid"
	"github.com/gorilla/mux"
)

type apiGroupResponse struct {
	API    string
	Total  int
	Entity []Group
}

// Companies list of companies
func Groups(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Println("GET Companies")

		var groups []Group
		Db.Find(&groups)

		var Response apiGroupResponse
		Response.Entity = groups
		Response.API = Version
		Response.Total = len(Response.Entity)

		addedrecordString, _ := json.Marshal(Response)

		w.WriteHeader(http.StatusOK)
		n, _ := fmt.Fprintf(w, string(addedrecordString))
		log.Println(n)

	default:
		n, _ := fmt.Fprintf(w, "Sorry, only GET method are supported.")
		log.Println(n)
	}
}

// GroupOptionGetPatchDelete crud for user object
func GroupOptionGetPatchDelete(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	fmt.Println(params)
	fmt.Printf("Req: %s %s\n", r.Host, r.URL.Path)

	switch r.Method {

	case "OPTIONS":
		log.Println("OPTIONS /group/")

	case "PATCH":
		log.Println("PATCH /group/" + params["id"])
		type patchCompanyObj struct {
			Name string
		}
		var ci patchCompanyObj
		err := json.NewDecoder(r.Body).Decode(&ci)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			n, _ := fmt.Fprintf(w, "{\"message\" : \"Unexpected Object Received\"}")
			fmt.Println(n)
			return
		}
		var record Group
		Db.Model(&record).Where("id = ?", params["id"]).Updates(User{Name: ci.Name})
		log.Println("PATCH /group/" + params["id"] + " DONE")

	case "DELETE":
		log.Println("DELETE /group/" + params["id"])
		var record Group
		Db.Where("id = ?", params["id"]).Delete(&record)
		log.Println("DELETE /group/" + params["id"] + " DONE")

	case "GET":
		log.Println("GET /group/" + params["id"])
		// id, _ := strconv.ParseUint(params["id"], 10, 64)
		var record Group
		Db.Where("id = ?", params["id"]).Last(&record)
		if record.Name != "" {
			addedrecordString, _ := json.Marshal(record)
			w.WriteHeader(http.StatusOK)
			n, _ := fmt.Fprintf(w, string(addedrecordString))
			log.Println(n)
		} else {
			w.WriteHeader(http.StatusNotFound)
			n, _ := fmt.Fprintf(w, "{\"message\":\"Group Not Found\"}")
			log.Println(n)
		}

	default:
		n, _ := fmt.Fprintf(w, "Sorry, only PATCH,DELETE,GET methods are supported.")
		log.Println(n)
	}
}

// GroupPostOptions post new user
func GroupPostOptions(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "OPTIONS":
		log.Println("OPTIONS /group")

	case "POST":
		log.Println("POST /group")

		type addCompanyObj struct {
			Name string
		}
		var ci addCompanyObj
		err := json.NewDecoder(r.Body).Decode(&ci)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			n, _ := fmt.Fprintf(w, "{\"message\" : \"Unexpected Object Received\"}")
			log.Println(n)
			return
		}

		cuserinterface := r.Context().Value("user")
		fmt.Println(cuserinterface)
		// cid, _ := strconv.ParseUint(*cuserinterface.(*string), 10, 64)
		// cid, _ := cuserinterface.(uuid.UUID)

		var group Group
		group.Name = ci.Name
		// group.CreatedBy = cid
		// group.UpdatedBy = cid

		dberrors := Db.Create(&group).GetErrors()
		if len(dberrors) != 0 {
			addedrecordString, _ := json.Marshal(dberrors)
			w.WriteHeader(http.StatusInternalServerError)
			n, _ := fmt.Fprintf(w, string(addedrecordString))
			log.Println(n)
		} else {
			addedrecordString, _ := json.Marshal(group)
			w.WriteHeader(http.StatusCreated)
			n, _ := fmt.Fprintf(w, string(addedrecordString))
			log.Println(n)
		}

	default:
		n, _ := fmt.Fprintf(w, "Sorry, only OPTIONS,POST methods are supported.")
		log.Println(n)
	}
}
