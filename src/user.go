package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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
	//case "OPTIONS":

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
		n, _ := fmt.Fprintf(w, string(addedrecordString))
		fmt.Println(n)
	default:
		n, _ := fmt.Fprintf(w, "Sorry, only GET method are supported.")
		fmt.Println(n)
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
		buf, _ := ioutil.ReadAll(r.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
		type patchType struct {
			PropertyType string
		}
		var incType patchType
		err := json.NewDecoder(rdr1).Decode(&incType)
		if err != nil || incType.PropertyType == "" {
			w.WriteHeader(http.StatusBadRequest)
			n, _ := fmt.Fprintf(w, "{\"message\" : \"Unexpected Object Received\"}")
			log.Println(n)
			return
		}

		var user User
		switch incType.PropertyType {
		case "STRING":
			type patchTextField struct {
				PropertyName  string
				PropertyValue string
			}
			var incomingData patchTextField
			err := json.NewDecoder(rdr2).Decode(&incomingData)
			if err != nil || incomingData.PropertyName == "" {
				w.WriteHeader(http.StatusBadRequest)
				n, _ := fmt.Fprintf(w, "{\"message\" : \"Unexpected Object Received\"}")
				log.Println(n)
				return
			}

			Db.Where("id = ?", params["id"]).Model(&user).Update(incomingData.PropertyName, incomingData.PropertyValue)

		case "BOOLEAN":
			type patchTextField struct {
				PropertyName  string
				PropertyValue bool
			}
			var incomingData patchTextField
			err := json.NewDecoder(rdr2).Decode(&incomingData)
			if err != nil || incomingData.PropertyName == "" {
				w.WriteHeader(http.StatusBadRequest)
				n, _ := fmt.Fprintf(w, "{\"message\" : \"Unexpected Object Received\"}")
				log.Println(n)
				return
			}

			Db.Where("id = ?", params["id"]).Model(&user).Update(incomingData.PropertyName, incomingData.PropertyValue)
		}

		log.Println("PATCH /user/" + params["id"] + " DONE")

		patchedRecordString, _ := json.Marshal(user)
		w.WriteHeader(http.StatusOK)
		n, _ := fmt.Fprintf(w, string(patchedRecordString))
		fmt.Println(n)

	case "DELETE":
		log.Println("DELETE /user/" + params["id"])
		Db.Where("id = ?", params["id"]).Delete(&User{})
		w.WriteHeader(http.StatusOK)
		n, _ := fmt.Fprintf(w, "")
		fmt.Println(n)

	case "GET":
		log.Println("GET /user/" + params["id"])
		// id, _ := strconv.ParseUint(params["id"], 10, 64)
		var user User
		Db.Where("id = ?", params["id"]).Last(&user)
		if user.Name != "" {
			addedrecordString, _ := json.Marshal(user)
			w.WriteHeader(http.StatusOK)
			n, _ := fmt.Fprintf(w, string(addedrecordString))
			fmt.Println(n)
		} else {
			w.WriteHeader(http.StatusNotFound)
			n, _ := fmt.Fprintf(w, "{\"message\":\"User Not Found\"}")
			fmt.Println(n)
		}

	default:
		n, _ := fmt.Fprintf(w, "Sorry, only PATCH,DELETE,GET methods are supported.")
		fmt.Println(n)
	}
}

func Invite(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	switch r.Method {
	case "GET":
		var user User
		Db.Where("id = ?", params["id"]).Find(&user)
		if user.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			msg := "Not Found"
			if DEV {
				msg = msg + "\nUser not found, no one associated with that id"
			}
			fmt.Println(msg)
			n, _ := fmt.Fprintf(w, "{\"message\" : \""+msg+"\"}")
			log.Println(n)
			return
		}
		user.SetPass = GetHash()
		Db.Save(&user)
		type inviteResponse struct {
			Link string
			Hash string
		}
		var invResp inviteResponse
		invResp.Hash = user.SetPass
		invResp.Link = os.Getenv("HOST") + "/password/" + user.SetPass

		addedRecordString, err := json.Marshal(invResp)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			n, _ := fmt.Fprintf(w, string(addedRecordString))
			fmt.Println(n)
			return
		}
	}
}

// SetPassword special method to set password
func SetPassword(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	switch r.Method {
	case "GET":

		if params["id"] == "" {
			w.WriteHeader(http.StatusBadRequest)
			msg := "Set Password Bad Request"
			if DEV {
				msg = msg + "\nFor Some reason Hash (id) is empty" + params["id"]
			}
			fmt.Println(msg)
			n, _ := fmt.Fprintf(w, "{\"message\" : \""+msg+"\"}")
			log.Println(n)
			return
		}

		var user User
		Db.Where("set_pass = ?", params["id"]).Find(&user)
		if user.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			msg := "Set Password Bad Request"
			if DEV {
				msg = msg + "\nThere is no user associated with incoming hash1 " + params["id"]
			}
			fmt.Println(msg)
			n, _ := fmt.Fprintf(w, "{\"message\" : \""+msg+"\"}")
			log.Println(n)
			return
		} else {
			user.SetPass = GetHash()
			Db.Save(&user)
		}
		tmpl := template.Must(template.ParseFiles("password.html"))

		type TemplateData struct {
			Host string
			Hash string
			Code string
		}
		var data TemplateData
		data.Host = os.Getenv("HOST")
		data.Hash = params["id"]
		data.Code = user.SetPass
		w.Header().Set("content-type", "content-type: text/html;")
		tmpl.Execute(w, data)

	case "POST":
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			msg := "Unexpected Object Received"
			if DEV {
				if r.FormValue("Code") == "" {
					msg = msg + "\nEmpty Code received"
				}
				if r.FormValue("pwd1") == "" {
					msg = msg + "\nEmpty Password received"
				}
				if r.FormValue("pwd2") == "" {
					msg = msg + "\nEmpty Password received"
				}
				if r.FormValue("pwd2") != r.FormValue("pwd1") {
					msg = msg + "\nPasswords doesn't match"
				}
			}
			fmt.Println(msg)
			n, _ := fmt.Fprintf(w, "{\"message\" : \""+msg+"\"}")
			log.Println(n)
			return
		}

		var user User
		Db.Where("set_pass = ?", r.FormValue("Code")).Find(&user)
		if user.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			msg := "Unexpected Object Received"
			if DEV {
				msg = msg + "\nThere is no user associated with that hash2 " + params["id"] + " in the database"
			}
			fmt.Println(msg)
			n, _ := fmt.Fprintf(w, "{\"message\" : \""+msg+"\"}")
			log.Println(n)
			return
		}
		user.Hash = HashAndSalt([]byte(r.FormValue("pwd1")))
		user.SetPass = ""
		user.Active = true
		Db.Save(&user)
		w.WriteHeader(http.StatusOK)
		n, _ := fmt.Fprintf(w, "")
		log.Println(n)
		return
	}
}

// UserPostOptions post new user
func UserPostOptions(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "OPTIONS":
		log.Println("OPTIONS /user")
	case "POST":
		log.Println("POST /user")

		type addUserObj struct {
			Name     string
			FullName string
			Email    string
		}
		var incomingData addUserObj
		err := json.NewDecoder(r.Body).Decode(&incomingData)
		if err != nil || incomingData.Name == "" || incomingData.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			n, _ := fmt.Fprintf(w, "{\"message\" : \"Unexpected Object Received\"}")
			log.Println(n)
			return
		}

		cuserinterface := r.Context().Value("user").(string)
		fmt.Println("cuserinterface")
		fmt.Println(cuserinterface)

		var user User
		user.Name = incomingData.Name
		user.FullName = incomingData.FullName
		user.Email = incomingData.Email
		user.CreatedBy = cuserinterface
		user.UpdatedBy = cuserinterface
		user.SetPass = GetHash()

		dberrors := Db.Create(&user).GetErrors()
		if len(dberrors) != 0 {
			addedrecordString, _ := json.Marshal(dberrors)
			w.WriteHeader(http.StatusInternalServerError)
			n, _ := fmt.Fprintf(w, string(addedrecordString))
			log.Println(n)
		} else {
			addedrecordString, _ := json.Marshal(user)
			w.WriteHeader(http.StatusCreated)
			n, _ := fmt.Fprintf(w, string(addedrecordString))
			log.Println(n)
		}

	default:
		n, _ := fmt.Fprintf(w, "Sorry, only OPTIONS,POST methods are supported.")
		fmt.Println(n)
	}
}
