package key

import (
	"GORM-practice-backend/app/modules/key"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"GORM-practice-backend/app/helpers"
	"GORM-practice-backend/app/models"
	"GORM-practice-backend/app/modules/auth"
	"github.com/gorilla/mux"
)

//TO DO IN KEY PACKAGE
//Update Key By ID (Input Using GORM, more to it later)
//Get Key By Filter(Name, Type, Platform, App Version, Tribe, Status)

// CreateKeyHandler create key
func (h *Handler) CreateKeyHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	message := JSONMessage{
		Status:  "Success",
		Message: "Created Key",
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][CreateKeyHandler][ReadBody]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when creating key"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	}

	key := models.Key{}
	if err = json.Unmarshal(body, &key); err != nil {
		fmt.Printf("[crud_key_handler.go][CreateKeyHandler][UnmarshalJSON]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when creating key"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	}
	uid, _, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"error UID extraction",
		}`), http.StatusInternalServerError)
	}

	key.UserID = uint(uid)

	if err = h.CreateKey(key); err != nil {
		fmt.Printf("[crud_key_handler.go][CreateKeyHandler][InsertTribe]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when creating key"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	}
	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	return
}

// DeleteKeyHandler delete key
func (h *Handler) DeleteKeyHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	message := JSONMessage{
		Status:  "Success",
		Message: "Deleted Key",
	}

	params := mux.Vars(r)
	targetUint, err := strconv.ParseUint(params["key_id"], 10, 32)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][DeleteKeyHandler][ParseUint]: %s", err)
		message.Status = "Failed"
		message.Message = "Error while deleting"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	if err = h.DeleteKey(uint(targetUint)); err != nil {
		fmt.Printf("[crud_key_handler.go][DeleteKeyHandler][DeleteTribe]: %s", err)
		message.Status = "Failed"
		message.Message = "Error while deleting"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	return
}

//GetKeyByID by user
func (h *Handler) GetKeyByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var key models.Key
	h.DB.Preload("Shares").First(&key, params["key_id"])
	json.NewEncoder(w).Encode(&key)
}

//UpdateKeyByID key
func (h *Handler) UpdateKeyByID(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	message := JSONMessage{
		Status:  "Success",
		Message: "Updated Key",
	}

	params := mux.Vars(r)
	var key models.Key
	h.DB.First(&key, params["key_id"])
	//read edit info
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][UpdateKeyByID][ReadBody]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when updating key"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	}

	updateKey := models.Key{}
	if err = json.Unmarshal(body, &updateKey); err != nil {
		fmt.Printf("[crud_key_handler.go][UpdateKeyByID][UnmarshalJSON]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when updating key"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	}

	//Checks things to be updated
	if updateKey.KeyName != "" {
		key.KeyName = updateKey.KeyName
	}
	if updateKey.KeyValue != "" {
		key.KeyValue = updateKey.KeyValue
	}
	if updateKey.KeyType != "" {
		key.KeyType = updateKey.KeyType
	}
	if updateKey.Description != "" {
		key.Description = updateKey.Description
	}
	if updateKey.Platform != "" {
		key.Platform = updateKey.Platform
	}
	if updateKey.ExpireDate.IsZero() {
		key.ExpireDate = updateKey.ExpireDate
	}
	if updateKey.UserID != 0 {
		key.UserID = updateKey.UserID
	}
	if updateKey.TribeID != 0 {
		key.TribeID = updateKey.TribeID
	}
	if updateKey.AppVersion != "" {
		key.AppVersion = updateKey.AppVersion
	}
	if updateKey.Status != "" {
		key.Status = updateKey.Status
	}
	h.DB.Save(&key)

	message = JSONMessage{
		Status:  "Success",
		Message: "Updated Key",
	}
	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	return
}

// GetKeysByUserID as said
func (h *Handler) GetKeysByUserID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var keys []models.Key

	uid, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"error UID extraction",
		}`), http.StatusInternalServerError)
		return
	}

	if role < 1 || string(uid) != params["user_id"] { // Get user own key
		helpers.RenderJSON(w, []byte(`
		{
			"message":"you are not authorized to request",
		}`), http.StatusUnauthorized)
		return
	}

	h.DB.Preload("Shares").Where("user_id = ?", params["user_id"]).Find(&keys)
	json.NewEncoder(w).Encode(&keys)
}

// GetKeysByTribeID as said
func (h *Handler) GetKeysByTribeID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var keys []models.Key

	h.DB.Preload("Shares").Where("tribe_id = ?", params["tribe_id"]).Find(&keys)
	json.NewEncoder(w).Encode(&keys)
}

func (h *Handler) ShareKey(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	message := JSONMessage{
		Status:  "Success",
		Message: "Shared Key Successfully",
	}

	//get tribe uint64
	params := mux.Vars(r)
	keyUint, err := strconv.ParseUint(params["tribe_id"], 10, 32)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][ShareKey][ParseUint]: %s", err)
		message.Status = "Failed"
		message.Message = "Error when trying to share key"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][ShareKey][ReadBody]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when trying to share key"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	}

	var assign Assign
	//read body, get user id
	if err = json.Unmarshal(body, &assign); err != nil {
		fmt.Printf("[crud_key_handler.go][ShareKey][UnmarshalJSON]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when trying to share key"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	}

	var tribe models.Tribe
	h.DB.First(&key, uint(keyUint))
	h.DB.Model(&key).Association("Shares").Append(models.KeyShares{UserID: assign.UID, KeyID: uint(KeyUint)})

	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	return
}
