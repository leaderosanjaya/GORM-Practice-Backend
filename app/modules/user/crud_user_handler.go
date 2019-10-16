package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/GORM-practice/app/helpers"
	"github.com/GORM-practice/app/models"
	"github.com/GORM-practice/app/modules/auth"

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

//CreateUserHandler create user
func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[CRUD User Read Body][User]: %s", err)
		helpers.SendError(w, "error insert user", http.StatusBadRequest)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Printf("[CRUD User Unmarshal JSON][User]: %s", err)
		helpers.SendError(w, "error insert user", http.StatusBadRequest)
		return
	}

	cred := Credential{}
	err = json.Unmarshal(body, &cred)
	if err != nil {
		fmt.Printf("[CRUD User Unmarshal JSON][Cred]: %s", err)
		helpers.SendError(w, "error insert user", http.StatusBadRequest)

		return
	}

	user.Role = 0
	user.Password = cred.Password

	user = models.User(user)
	if err = h.InsertUser(user); err != nil {
		fmt.Printf("[CRUD User Insert User][User]: %s", err)
		helpers.SendError(w, "error insert user", http.StatusInternalServerError)
		return
	}

	helpers.SendOK(w, "user registered")
	return
}

// DeleteUserHandler delete user handler
func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	targetUint, err := strconv.ParseUint(params["user_id"], 10, 32)
	if err != nil {
		fmt.Printf("[crud_user_handler.go][DeleteUserHandler][ParseUint]: %s", err)
		helpers.SendError(w, "error delete user", http.StatusBadRequest)
		return
	}

	_, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.SendError(w, "error delete user", http.StatusBadRequest)
		return
	}

	if role != 1 {
		helpers.SendError(w, "super admin access only", http.StatusForbidden)
		return
	}

	if err = h.DeleteUser(uint(targetUint)); err != nil {
		fmt.Printf("[CRUD User Delete User][User]: %s", err)
		helpers.SendError(w, "error delete user", http.StatusInternalServerError)
		return
	}

	helpers.SendOK(w, "user deleted")
	return
}

//GetUserByID get user by id
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user models.User
	h.DB.Preload("Keys").Preload("Tribes").Preload("SharedKeys").First(&user, params["user_id"])
	write, _ := json.Marshal(&user)
	helpers.RenderJSON(w, write, http.StatusOK)
}

// IMPROVE
func (h *Handler) UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user models.User
	h.DB.First(&user, params["user_id"])

	uid, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.SendError(w, "error UID extraction", http.StatusInternalServerError)
		return
	}

	if role < 1 && uint64(user.ID) != uid { // Get user own key
		helpers.SendError(w, "You are not authorized for this request", http.StatusUnauthorized)
		return
	}

	//read edit info
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_user_handler.go][UpdateUserByID][ReadBody]: %s\n", err)
		helpers.SendError(w, "Error when updating user", http.StatusBadRequest)
		return
	}

	updateUser := models.User{}
	if err = json.Unmarshal(body, &updateUser); err != nil {
		fmt.Printf("[crud_user_handler.go][UpdateUserByID][UnmarshalJSON]: %s\n", err)
		helpers.SendError(w, "Error when updating user", http.StatusBadRequest)
		return
	}

	UpdateValue(&updateUser, &user)
	h.DB.Save(&user)
	helpers.SendOK(w, "Updated user")
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Get User ID
	_, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.SendError(w, "error UID extraction", http.StatusInternalServerError)
		return
	}
	if role < 1 {
		helpers.SendError(w, "Request denied, superadmin only", http.StatusUnauthorized)
		return
	}

	var users []models.User
	h.DB.Preload("Keys").Preload("Tribes").Preload("SharedKeys").Where("role = ?", "0").Order("user_id desc").Find(&users)
	write, _ := json.Marshal(&users)
	helpers.RenderJSON(w, write, http.StatusOK)
}

// GetTribeByUser get users tribe
func (h *Handler) GetTribeByUser(w http.ResponseWriter, r *http.Request) {
	// Get User ID
	uid, _, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.SendError(w, "error uid extraction", http.StatusInternalServerError)
		return
	}

	var tribeAss []models.TribeAssign
	var gotTribe = false
	if row := h.DB.Table("tribe_assigns").Find(&tribeAss, uid).RowsAffected; row != 0 {
		gotTribe = true
	}

	resp := map[string]interface{}{"status": gotTribe}
	resp["tribeAssign"] = tribeAss
	write, _ := json.Marshal(resp)
	helpers.RenderJSON(w, write, http.StatusOK)
}

// GetTribeByUserID get users tribe
func (h *Handler) GetTribeByUserID(w http.ResponseWriter, r *http.Request) {
	// Get User ID
	_, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.SendError(w, "error uid extraction", http.StatusInternalServerError)
		return
	}
	if role < 1 {
		helpers.SendError(w, "super admin access only", http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	var tribeAss []models.TribeAssign
	var gotTribe = false

	if row := h.DB.Table("tribe_assigns").Find(&tribeAss, params["user_id"]).RowsAffected; row != 0 {
		gotTribe = true
	}

	resp := map[string]interface{}{"status": gotTribe}
	resp["tribeAssign"] = tribeAss
	write, _ := json.Marshal(resp)
	helpers.RenderJSON(w, write, http.StatusOK)
}
