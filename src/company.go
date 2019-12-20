package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// "github.com/google/uuid"
	"github.com/gorilla/mux"
)

type apiCompanyResponse struct {
	API    string
	Total  int
	Entity []Company
}

// Companies list of companies
func Companies(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Println("GET Companies")

		var companies []Company
		Db.Find(&companies)

		var Response apiCompanyResponse
		Response.Entity = companies
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

// CompanyOptionGetPatchDelete crud for user object
func CompanyOptionGetPatchDelete(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	fmt.Println(params)
	fmt.Printf("Req: %s %s\n", r.Host, r.URL.Path)

	switch r.Method {

	case "OPTIONS":
		log.Println("OPTIONS /company/")

	case "PATCH":
		log.Println("PATCH /company/" + params["id"])
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
		var record Company
		Db.Model(&record).Where("id = ?", params["id"]).Updates(User{Name: ci.Name})
		log.Println("PATCH /company/" + params["id"] + " DONE")

	case "DELETE":
		log.Println("DELETE /company/" + params["id"])
		var record Company
		Db.Where("id = ?", params["id"]).Delete(&record)
		log.Println("DELETE /company/" + params["id"] + " DONE")

	case "GET":
		log.Println("GET /company/" + params["id"])
		// id, _ := strconv.ParseUint(params["id"], 10, 64)
		var record Company
		Db.Where("id = ?", params["id"]).Last(&record)
		if record.Name != "" {
			addedrecordString, _ := json.Marshal(record)
			w.WriteHeader(http.StatusOK)
			n, _ := fmt.Fprintf(w, string(addedrecordString))
			log.Println(n)
		} else {
			w.WriteHeader(http.StatusNotFound)
			n, _ := fmt.Fprintf(w, "{\"message\":\"Company Not Found\"}")
			log.Println(n)
		}

	default:
		n, _ := fmt.Fprintf(w, "Sorry, only PATCH,DELETE,GET methods are supported.")
		log.Println(n)
	}
}

// CompanyPostOptions post new user
func CompanyPostOptions(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "OPTIONS":
		log.Println("OPTIONS /company")

	case "POST":
		log.Println("POST /company")

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

		var company Company
		company.Name = ci.Name
		// company.CreatedBy = cid
		// company.UpdatedBy = cid

		dberrors := Db.Create(&company).GetErrors()
		if len(dberrors) != 0 {
			addedrecordString, _ := json.Marshal(dberrors)
			w.WriteHeader(http.StatusInternalServerError)
			n, _ := fmt.Fprintf(w, string(addedrecordString))
			log.Println(n)
		} else {
			addedrecordString, _ := json.Marshal(company)
			w.WriteHeader(http.StatusCreated)
			n, _ := fmt.Fprintf(w, string(addedrecordString))
			log.Println(n)
		}

	default:
		n, _ := fmt.Fprintf(w, "Sorry, only OPTIONS,POST methods are supported.")
		log.Println(n)
	}
}
