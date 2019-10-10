package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"GORM-practice-backend/app/models"
	"GORM-practice-backend/app/modules/auth"
	"GORM-practice-backend/app/modules/key"
	"GORM-practice-backend/app/modules/remote-config"
	"GORM-practice-backend/app/modules/tribe"
	"GORM-practice-backend/app/modules/user"
	"GORM-practice-backend/config"
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
	remoteConfigHandler := new(remoteconfig.Handler)
	remoteConfigHandler.Init()
	//Pass DB to handler
	userHandler.DB = db
	tribeHandler.DB = db
	keyHandler.DB = db
	authHandler.DB = db
	remoteConfigHandler.DB = db
	keyHandler.PushRemoteConfig = remoteConfigHandler.PublishConfig

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
	s.HandleFunc("/api/users/{user_id:[0-9]+}", userHandler.GetUserByID).Methods("GET")
	//Get user keys by ID
	//IMPLEMENT FILTER SOON.
	s.HandleFunc("/api/users/{user_id:[0-9]+}/keys", keyHandler.GetKeysByUserID).Methods("GET")
	s.HandleFunc("/api/tribes/{tribe_id:[0-9]+}/keys", keyHandler.GetKeysByTribeID).Methods("GET")
	//Get shared keys by ID

	//Get user by filter
	// router.HandleFunc("/api/users", userHandler.GetUsers).Methods("GET")
	// router.HandleFunc("/api/user/{user_id:[0-9]+}/tribes")

	//Create and delete tribe
	s.HandleFunc("/api/tribes", tribeHandler.CreateTribeHandler).Methods("POST")
	s.HandleFunc("/api/tribes/{tribe_id:[0-9]+}", tribeHandler.DeleteTribeHandler).Methods("DELETE")

	//Assign user to tribe
	s.HandleFunc("/api/tribes/{tribe_id:[0-9]+}/members", tribeHandler.AssignUser).Methods("POST")
	//Remove user from tribe
	s.HandleFunc("/api/tribes/{tribe_id:[0-9]+}/members", tribeHandler.RemoveAssign).Methods("DELETE")
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
	s.HandleFunc("/api/keys/{key_id:[0-9]+}", keyHandler.GetKeyByID).Methods("GET")
	//Update Key by ID
	s.HandleFunc("/api/keys/{key_id:[0-9]+}", keyHandler.UpdateKeyByID).Methods("PUT")

	//Assign Key Share
	s.HandleFunc("/api/keys/{key_id:[0-9]+}/shares", keyHandler.ShareKey).Methods("POST")
	//Remove Key Share
	s.HandleFunc("/api/keys/{key_id:[0-9]+}/shares", keyHandler.RevokeShare).Methods("DELETE")

	//Get keys by filter
	// router.HandleFunc("/api/keys/").Methods("GET")

	//Get key user functions
	//Read Key list given User, use Limit and Page

	//Delete Key by Name

	//Update Key by Name, given new value
	//change it to execute from createkey
	// err = remoteConfigHandler.PublishConfig()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	port := os.Getenv("PORT") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8080" //localhost
	}
	fmt.Printf("[%s] Listening on Port 8080\n", time.Now())
	log.Fatal(http.ListenAndServe(":"+port, gorillaHandler.CORS(headers, methods, origins)(router)))
}
