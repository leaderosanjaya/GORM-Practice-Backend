package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/GORM-practice/app/models"
	"github.com/GORM-practice/app/modules/auth"
	"github.com/GORM-practice/app/modules/key"
	"github.com/GORM-practice/app/modules/tribe"
	"github.com/GORM-practice/app/modules/user"
	"github.com/GORM-practice/config"
	gorillaHandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API Calls for GORM Practice.")
}

func main() {
	//Connect DB
	db, err := config.ConnectDB()
	if err != nil {
		fmt.Printf("[main.go][ConnectDB]: %s\n", err)
	}

	defer db.Close()
	// new handler
	userHandler := new(user.Handler)
	tribeHandler := new(tribe.Handler)
	keyHandler := new(key.Handler)
	authHandler := new(auth.Handler)

	//Pass DB to handler
	userHandler.DB = db
	tribeHandler.DB = db
	keyHandler.DB = db
	authHandler.DB = db

	//Update schema to models.go
	db.AutoMigrate(&models.User{}, &models.Tribe{}, &models.Key{}, &models.KeyShares{}, &models.TribeAssign{})
	db.Model(&models.KeyShares{}).AddForeignKey("user_id", "users(user_id)", "CASCADE", "CASCADE")
	db.Model(&models.KeyShares{}).AddForeignKey("key_id", "keys(key_id)", "CASCADE", "CASCADE")

	db.Model(&models.TribeAssign{}).AddForeignKey("user_id", "users(user_id)", "CASCADE", "CASCADE")
	db.Model(&models.TribeAssign{}).AddForeignKey("tribe_id", "tribes(tribe_id)", "CASCADE", "CASCADE")

	//New Router
	router := mux.NewRouter()
	headers := gorillaHandler.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := gorillaHandler.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := gorillaHandler.AllowedOrigins([]string{"*"})

	s := router.PathPrefix("/auth").Subrouter()
	s.Use(auth.JwtVerify)

	router.HandleFunc("/api", index).Methods("GET")

	//Login user
	router.HandleFunc("/api/login", authHandler.Login).Methods("POST")

	//Create and delete user
	router.HandleFunc("/api/users", userHandler.CreateUserHandler).Methods("POST")
	s.HandleFunc("/api/users/{user_id:[0-9]+}", userHandler.DeleteUserHandler).Methods("DELETE")
	//Get user By ID
	router.HandleFunc("/api/users/{user_id:[0-9]+}", userHandler.GetUserByID).Methods("GET")
	//Get user keys by ID
	//IMPLEMENT FILTER SOON.
	router.HandleFunc("/api/users/{user_id:[0-9]+}/keys", keyHandler.GetKeysByUserID).Methods("GET")
	router.HandleFunc("/api/tribes/{tribe_id:[0-9]+}/keys", keyHandler.GetKeysByTribeID).Methods("GET")
	//Get shared keys by ID

	//Get user by filter
	// router.HandleFunc("/api/users", userHandler.GetUsers).Methods("GET")
	// router.HandleFunc("/api/user/{user_id:[0-9]+}/tribes")

	//Create and delete tribe
	router.HandleFunc("/api/tribes", tribeHandler.CreateTribeHandler).Methods("POST")
	router.HandleFunc("/api/tribes/{tribe_id:[0-9]+}", tribeHandler.DeleteTribeHandler).Methods("DELETE")

	//Assign user to tribe
	router.HandleFunc("/api/tribes/{tribe_id:[0-9]+}/members", tribeHandler.AssignUser).Methods("POST")
	//Remove user from tribe
	router.HandleFunc("/api/tribes/{tribe_id:[0-9]+}/members", tribeHandler.RemoveAssign).Methods("DELETE")
	//get tribe by id
	router.HandleFunc("/api/tribes/{tribe_id:[0-9]+}", tribeHandler.GetTribeByID).Methods("GET")
	//Get tribe keys
	// router.HandleFunc("/api/tribe/{tribe_id:[0-9]+}/keys").Methods("GET")
	//Get tribe users
	// router.HandleFunc("/api/tribe/{tribe_id:[0-9]+}/users").Methods("GET")
	//Create and delete Key
	s.HandleFunc("/api/keys", keyHandler.CreateKeyHandler).Methods("POST")
	s.HandleFunc("/api/keys/{key_id:[0-9]+}", keyHandler.DeleteKeyHandler).Methods("DELETE")
	//Get key by ID
	router.HandleFunc("/api/keys/{key_id:[0-9]+}", keyHandler.GetKeyByID).Methods("GET")
	//Update Key by ID
	router.HandleFunc("/api/keys/{key_id:[0-9]+}", keyHandler.UpdateKeyByID).Methods("PUT")

	//Assign Key Share
	router.HandleFunc("/api/keys/{key_id:[0-9]+}/shares", keyHandler.ShareKey).Methods("POST")
	//Remove Key Share
	router.HandleFunc("/api/keys/{key_id:[0-9]+}/shares", keyHandler.RevokeShare).Methods("DELETE")
	//Get keys by filter
	// router.HandleFunc("/api/keys/").Methods("GET")

	//Get key user functions
	//Read Key list given User, use Limit and Page

	//Delete Key by Name

	//Update Key by Name, given new value

	fmt.Printf("[%s] Listening on Port 8080", time.Now())
	log.Fatal(http.ListenAndServe(":8080", gorillaHandler.CORS(headers, methods, origins)(router)))
}
