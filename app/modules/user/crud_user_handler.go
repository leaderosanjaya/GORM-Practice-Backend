package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/GORM-practice/app/helpers"
	"github.com/GORM-practice/app/models"
	"github.com/gorilla/mux"
)

// REFACTOR CODE
// Error convention: [FileName][FunctionName][Process]:
// Doing it later as I am too lazy

// func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
// 	status := http.StatusOK
// 	message := []byte("")

// 	users, err := h.GetUsers()

// }

func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	message := JSONMessage{
		Status:  "Success",
		Message: "Registered User",
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[CRUD User Read Body][User]: %s", err)
		message.Status = "Failed"
		message.Message = "Error while registering"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Printf("[CRUD User Unmarshal JSON][User]: %s", err)
		message.Status = "Failed"
		message.Message = "Error while registering"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	if err = h.InsertUser(user); err != nil {
		fmt.Printf("[CRUD User Insert User][User]: %s", err)
		message.Status = "Failed"
		message.Message = "Error while registering"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	return
}

func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	message := JSONMessage{
		Status:  "Success",
		Message: "Deleted User",
	}

	params := mux.Vars(r)

	if err = h.DeleteUser(params["user_id"]); err != nil {
		fmt.Printf("[CRUD User Insert User][User]: %s", err)
		message.Status = "Failed"
		message.Message = "Error while deleting"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	return
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user models.User
	h.DB.Preload("Keys").Preload("Tribes").Preload("SharedKeys").First(&user, params["user_id"])
	json.NewEncoder(w).Encode(&user)
}

// func (h *Handler) GetUserKeyByID(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	var keys
// }
