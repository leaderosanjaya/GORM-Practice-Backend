package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GORM-practice/app/helpers"
	"github.com/GORM-practice/app/models"
	"github.com/GORM-practice/app/modules/auth"
	"github.com/GORM-practice/app/modules/key"
	"github.com/GORM-practice/app/modules/remote-config"
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
	// IMPROVE: CLEAN THIS UP
	// new handler
	userHandler := new(user.Handler)
	tribeHandler := new(tribe.Handler)
	keyHandler := new(key.Handler)
	authHandler := new(auth.Handler)
	remoteConfigHandler := new(remoteconfig.Handler)

	//Pass DB to handler
	userHandler.DB = db
	tribeHandler.DB = db
	keyHandler.DB = db
	authHandler.DB = db
	remoteConfigHandler.DB = db
	keyHandler.PushRemoteConfig = remoteConfigHandler.PublishConfig

	// IMPROVE PUT IN SINGLE FUNCTION
	//Update schema to models.go
	db.AutoMigrate(&models.User{}, &models.Tribe{}, &models.Key{}, &models.KeyShares{}, &models.TribeAssign{}, &models.TribeLeadAssign{})
	db.Model(&models.KeyShares{}).AddForeignKey("user_id", "users(user_id)", "CASCADE", "CASCADE")
	db.Model(&models.KeyShares{}).AddForeignKey("key_id", "keys(key_id)", "CASCADE", "CASCADE")

	db.Model(&models.TribeAssign{}).AddForeignKey("user_id", "users(user_id)", "CASCADE", "CASCADE")
	db.Model(&models.TribeAssign{}).AddForeignKey("tribe_id", "tribes(tribe_id)", "CASCADE", "CASCADE")

	remoteConfigHandler.Init()

	//New Router
	router := mux.NewRouter()

	headers := gorillaHandler.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := gorillaHandler.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := gorillaHandler.AllowedOrigins([]string{"*"})

	s := router.PathPrefix("/auth").Subrouter()
	s.Use(auth.JwtVerify)

	// route to test if API is alive or not
	router.HandleFunc("/api", index).Methods("GET")

	//Login user
	router.HandleFunc("/api/login", authHandler.Login).Methods("POST")

	//Validate token user has
	s.HandleFunc("/api/user/valid", authHandler.ValidateToken).Methods("GET")

	//Create User
	router.HandleFunc("/api/users", userHandler.CreateUserHandler).Methods("POST")

	// IMPROVE: Get users by filter
	// router.HandleFunc("/api/users", userHandler.GetUsers).Methods("GET")

	// Delete User
	s.HandleFunc("/api/users/{user_id:[0-9]+}", userHandler.DeleteUserHandler).Methods("DELETE")

	//Get All User
	s.HandleFunc("/api/users", userHandler.GetAllUsers).Methods("GET")

	//Get user By ID
	s.HandleFunc("/api/users/{user_id:[0-9]+}", userHandler.GetUserByID).Methods("GET")

	//Update User
	s.HandleFunc("/api/users/{user_id:[0-9]+}", userHandler.UpdateUserByID).Methods("PUT")

	//Get user keys by ID
	// TODO implement filter
	s.HandleFunc("/api/users/{user_id:[0-9]+}/keys", keyHandler.GetKeysByUserID).Methods("GET")

	// IMPROVE: Get user affiliated tribes
	// router.HandleFunc("/api/user/{user_id:[0-9]+}/tribes")

	// Get keys from tribe ID
	// TODO implement filter
	s.HandleFunc("/api/tribes/{tribe_id:[0-9]+}/keys", keyHandler.GetKeysByTribeID).Methods("GET")

	//Create Tribe
	s.HandleFunc("/api/tribes", tribeHandler.CreateTribeHandler).Methods("POST")
	//Delete Tribe
	s.HandleFunc("/api/tribes/{tribe_id:[0-9]+}", tribeHandler.DeleteTribeHandler).Methods("DELETE")

	//Update Tribe
	s.HandleFunc("/api/tribes/{tribe_id:[0-9]+}", tribeHandler.UpdateTribeByID).Methods("PUT")

	// TODO: implement this
	// Get Tribes
	// s.HandleFunc("/api/tribes", tribeHandler.CreateTribeHandler).Methods("POST")

	// TODO ADD IN DOCS
	// Assign Lead
	s.HandleFunc("/api/tribes/{tribe_id:[0-9]+}/leads", tribeHandler.AddTribeLead).Methods("POST")
	// Remove Lead
	s.HandleFunc("/api/tribes/{tribe_id:[0-9]+}/leads", tribeHandler.RemoveTribeLead).Methods("DELETE")

	//Assign user to tribe
	s.HandleFunc("/api/tribes/{tribe_id:[0-9]+}/members", tribeHandler.AssignUser).Methods("POST")
	//Remove user from tribe
	s.HandleFunc("/api/tribes/{tribe_id:[0-9]+}/members", tribeHandler.RemoveAssign).Methods("DELETE")
	//get tribe by id
	s.HandleFunc("/api/tribes/{tribe_id:[0-9]+}", tribeHandler.GetTribeByID).Methods("GET")
	s.HandleFunc("/api/tribes/user", tribeHandler.GetTribeByUserID).Methods("GET")

	// TODO Get tribe keys
	// router.HandleFunc("/api/tribe/{tribe_id:[0-9]+}/keys").Methods("GET")
	// TODO Get tribe users
	// router.HandleFunc("/api/tribe/{tribe_id:[0-9]+}/users").Methods("GET")

	//Get All keys
	// TODO: ADD FILTER
	// TODO: FILTER BY tribe, version, key_type, platform, status
	s.HandleFunc("/api/keys", keyHandler.GetKeysHandler).Methods("GET")

	//Create New key
	s.HandleFunc("/api/keys", keyHandler.CreateKeyHandler).Methods("POST")

	//Delete Key by ID
	s.HandleFunc("/api/keys/{key_id:[0-9]+}", keyHandler.DeleteKeyHandler).Methods("DELETE")

	//Get Key by ID
	s.HandleFunc("/api/keys/{key_id:[0-9]+}", keyHandler.GetKeyByID).Methods("GET")

	//Update Key by ID
	s.HandleFunc("/api/keys/{key_id:[0-9]+}", keyHandler.UpdateKeyByID).Methods("PUT")

	//Assign Key Share
	s.HandleFunc("/api/keys/{key_id:[0-9]+}/shares", keyHandler.ShareKey).Methods("POST")

	//Remove Key Share
	s.HandleFunc("/api/keys/{key_id:[0-9]+}/shares", keyHandler.RevokeShare).Methods("DELETE")

	router.Use(helpers.LoggingMiddleware)
	port := os.Getenv("PORT") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8080" //localhost
	}
	fmt.Printf("[%s] Listening on Port 8080\n", time.Now())
	log.Fatal(http.ListenAndServe(":"+port, gorillaHandler.CORS(headers, methods, origins)(router)))
}
