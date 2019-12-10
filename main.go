package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	// "strings"

	"./src"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
)

func main() {



	
	err := godotenv.Load("parameters.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	src.AdmPass = os.Getenv("ADMP")




	src.Db, src.Err = gorm.Open("sqlite3", "w2g.db")
	if src.Err != nil {
		panic("failed to connect database")
	}
	defer src.Db.Close()
	src.Db.LogMode(src.DbLogMode)

	f, err := os.OpenFile("w2g.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	databasePrepare()

	handleHTTP()

	// Migrate the schema
	// src.Db.AutoMigrate(&boardTable{})

	// var id uint
	// row := src.Db.Table("roles").Where("role = ?", "Role 1").Select("id").Row()
	// row.Scan(&id)

	// var role src.Role
	// src.Db.Where("role = ?", "Role 1").Find(&role)

	// var users []src.User
	// // user.UserName = "User3"
	// // user.FullName = "Fill Name 3"
	// // user.RoleID = id

	// // src.Db.Create(&user)

	// src.Db.Model(&role).Related(&users)
	// log.Println(users)

}

func handleHTTP() {

	r := mux.NewRouter()
	r.Use(authMiddleware)
	r.Use(headerMiddleware)
	r.HandleFunc("/token", src.GetToken)
	r.HandleFunc("/refreshtoken", src.GetRefreshToken)
	r.HandleFunc("/users", src.Users)
	r.HandleFunc("/user/{id}", src.UserCRUD)

	fmt.Printf("Starting Server to HANDLE w2g.tech back end\nPort : " + src.Port + "\nAPI revision " + src.Version + "\n\n")
	if err := http.ListenAndServe(":"+src.Port, r); err != nil {
		log.Fatal(err)
	}
}

func databasePrepare() {

	src.Db.AutoMigrate(&src.User{})
	src.Db.AutoMigrate(&src.Role{})
	src.Db.AutoMigrate(&src.RolePermission{})
	src.Db.AutoMigrate(&src.UserPermission{})
	src.Db.AutoMigrate(&src.Token{})
	src.Db.AutoMigrate(&src.RefreshToken{})

	// type superAdminRolesType []src.Role
	var superAdminRoles []src.Role
	src.Db.Where("role = ?", "Super Administrator").Find(&superAdminRoles)
	var superAdminRoleID uint
	var superAdminRole src.Role
	if len(superAdminRoles) == 0 {
		superAdminRole.Role = "Super Administrator"
		superAdminRole.CreatedBy = 1
		src.Db.Create(&superAdminRole)
		superAdminRoleID = superAdminRole.ID
	} else {
		superAdminRoleID = superAdminRoles[0].ID
	}

	var user src.User
	src.Db.Where("user_name = ?", "w2g-admin").Last(&user)
	if user.UserName != "w2g-admin" {
		user.UserName = "w2g-admin"
		user.FullName = "Way2Go Super Admin"
		user.Email = "artem.kryhin@gmail.com"
		user.RoleID = superAdminRoleID
		user.Enabled = true
		user.CreatedBy = 1
		user.Hash = src.HashAndSalt([]byte(src.AdmPass))
		src.Db.Create(&user)

	} else {
		fmt.Println("Super Admin is On Board")
	}

	src.SetAllPermissionsToByRoleIDifNotExists(user.RoleID, "/user", true)
	src.SetAllPermissionsToByRoleIDifNotExists(user.RoleID, "/users", true)

	src.SetAllPermissionsToByRoleIDifNotExists(user.RoleID, "/role", true)
	src.SetAllPermissionsToByRoleIDifNotExists(user.RoleID, "/roles", true)


}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("URI " + r.RequestURI)
		if r.RequestURI == "/token" {
			next.ServeHTTP(w, r)
		} else {

			Authorization := r.Header.Get("Authorization")
			isUser, cUser := src.IsLegalUser(Authorization)

			if isUser {

				route := r.URL.Path
				method := r.Method

				if !src.IfUserHasPermission(cUser, src.GetRoute(route), method) {
					w.WriteHeader(http.StatusForbidden)
					fmt.Fprintf(w, "{\"message\":\"Access Denided\"}")
				}

				ctx := context.WithValue(r.Context(), "user", cUser)
				r = r.WithContext(ctx)

				next.ServeHTTP(w, r)

			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "{\"message\":\"User Not Found\"}")
			}
		}

	})
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		src.FillAnswerHeader(w)
		src.OptionsAnswer(w)

		// Do stuff here
		// log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
