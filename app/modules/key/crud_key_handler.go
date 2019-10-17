package key

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/GORM-practice/app/helpers"
	"github.com/GORM-practice/app/models"
	"github.com/GORM-practice/app/modules/auth"
	"github.com/gorilla/mux"
)

// TODO IN KEY PACKAGE
// TODO Get Key By Filter(Name, Type, Platform, App Version, Tribe, Status)

// CreateKeyHandler create key
func (h *Handler) CreateKeyHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][CreateKeyHandler][ReadBody]: %s\n", err)
		helpers.SendError(w, "Error when creating key", http.StatusBadRequest)
		return
	}

	key := models.Key{}
	if err = json.Unmarshal(body, &key); err != nil {
		fmt.Printf("[crud_key_handler.go][CreateKeyHandler][UnmarshalJSON]: %s\n", err)
		helpers.SendError(w, "Error when creating key", http.StatusBadRequest)
		return
	}

	// Check if key name contains spaces :: Firebase will return 400
	if strings.ContainsAny(key.KeyName, " ") {
		fmt.Printf("[crud_key_handler.go][CreateKeyHandler][UnmarshalJSON]: key name contains spaces\n")
		helpers.SendError(w, "Error when creating key", http.StatusBadRequest)
		return
	}

	// Get user ID from token
	uid, _, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.SendError(w, "error UID extraction", http.StatusInternalServerError)
		return
	}

	key.UserID = uint(uid)

	// var user models.TribeAssign
	// h.DB.Where("user_id = ?", uint(uid)).First(&user)
	// key.TribeID = user.TribeID

	if err = h.CreateKey(key); err != nil {
		fmt.Printf("[crud_key_handler.go][CreateKeyHandler][CreateKey]: %s\n", err)
		helpers.SendError(w, "Error when creating key", http.StatusBadRequest)
		return
	}

	err = h.PushRemoteConfig()
	if err != nil {
		fmt.Printf("[crud_key_handler.go][CreateKeyHandler][PushRemoteConfig]: %s\n", err)
		return
	}
	helpers.SendOK(w, "Key created Successfully")
}

// DeleteKeyHandler delete key
func (h *Handler) DeleteKeyHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	targetUint, err := strconv.ParseUint(params["key_id"], 10, 32)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][DeleteKeyHandler][ParseUint]: %s\n", err)
		helpers.SendError(w, "Error deleting key", http.StatusBadRequest)
		return
	}

	// Get user ID from token
	uid, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.SendError(w, "error UID extraction", http.StatusInternalServerError)
		return
	}

	var key models.Key
	row := h.DB.First(&key, params["key_id"]).RowsAffected
	if row == 0 {
		helpers.SendError(w, "Key does not exist", http.StatusBadRequest)
		return
	}

	if role < 1 && uint64(key.UserID) != uid {
		helpers.SendError(w, "Failed to delete, you are not the owner of this key", http.StatusUnauthorized)
		return
	}

	if err = h.DeleteKey(uint(targetUint)); err != nil {
		fmt.Printf("[crud_key_handler.go][DeleteKeyHandler][DeleteTribe]: %s\n", err)
		helpers.SendError(w, "Error deleting key", http.StatusInternalServerError)
		return
	}

	err = h.PushRemoteConfig()
	if err != nil {
		fmt.Printf("[crud_key_handler.go][DeleteKeyHandler][PushRemoteConfig]: %s\n", err)
		return
	}
	helpers.SendOK(w, "Key deleted successfully")
}

//GetKeyByID by user
func (h *Handler) GetKeyByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var key models.Key

	// Get user ID from token
	uid, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.SendError(w, "error UID extraction", http.StatusInternalServerError)
		return
	}

	row := h.DB.Preload("Shares").Preload("Conditions").First(&key, params["key_id"]).RowsAffected
	if row == 0 {
		helpers.SendError(w, "Key does not exist", http.StatusBadRequest)
		return
	}

	if key.UserID != uint(uid) && role < 1 {
		helpers.SendError(w, "You are not the owner of this key", http.StatusUnauthorized)
		return
	}

	write, _ := json.Marshal(&key)
	helpers.RenderJSON(w, write, http.StatusOK)
}

//UpdateKeyByID key
func (h *Handler) UpdateKeyByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var key models.Key
	h.DB.First(&key, params["key_id"])

	uid, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.SendError(w, "error UID extraction", http.StatusInternalServerError)
		return
	}

	if role < 1 && uint64(key.UserID) != uid { // Get user own key
		helpers.SendError(w, "You are not authorized for this request", http.StatusUnauthorized)
		return
	}

	//read edit info
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][UpdateKeyByID][ReadBody]: %s\n", err)
		helpers.SendError(w, "Error when updating key", http.StatusBadRequest)
		return
	}

	updateKey := models.Key{}
	if err = json.Unmarshal(body, &updateKey); err != nil {
		fmt.Printf("[crud_key_handler.go][UpdateKeyByID][UnmarshalJSON]: %s\n", err)
		helpers.SendError(w, "Error when updating key", http.StatusBadRequest)
		return
	}

	updateValue(&updateKey, &key)
	h.DB.Save(&key)

	err = h.PushRemoteConfig()
	if err != nil {
		fmt.Printf("[crud_key_handler.go][UpdateKeyByID][PushRemoteConfig]: %s\n", err)
		return
	}

	helpers.SendOK(w, "Updated key")
}

// GetKeysByUserID as said
func (h *Handler) GetKeysByUserID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var keys []models.Key

	uid, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.SendError(w, "error UID extraction", http.StatusInternalServerError)
		return
	}

	if paramUserID, _ := strconv.ParseUint(fmt.Sprintf("%s", params["user_id"]), 10, 32); paramUserID != uid && role < 1 { // Get user own key
		helpers.SendError(w, "You are not authorized for this request", http.StatusUnauthorized)
		return
	}

	h.DB.Preload("Shares").Preload("Conditions").Where("user_id = ?", params["user_id"]).Find(&keys)

	write, _ := json.Marshal(&keys)
	helpers.RenderJSON(w, write, http.StatusOK)
}

// GetKeysByTribeID as said
func (h *Handler) GetKeysByTribeID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var keys []models.Key

	uid, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.SendError(w, "error UID extraction", http.StatusInternalServerError)
		return
	}

	paramTribeID, _ := strconv.ParseUint(fmt.Sprintf("%s", params["tribe_id"]), 10, 32)

	if role < 1 {
		var userTribes []models.TribeAssign
		if row := h.DB.Where("user_id = ?", uint(uid)).Find(&userTribes); row.RowsAffected == 0 {
			helpers.SendError(w, "User does not exist", http.StatusBadRequest)
			return
		}
		if len(userTribes) == 0 {
			helpers.SendError(w, "User not in tribe, request denied", http.StatusUnauthorized)
			return
		}

		var ok = false
		for _, tribe := range userTribes {
			if uint64(tribe.TribeID) == paramTribeID {
				ok = true
				break
			}
		}
		if !ok {
			return
		}
	}

	h.DB.Preload("Shares").Preload("Conditions").Where("tribe_id = ?", params["tribe_id"]).Find(&keys)
	write, _ := json.Marshal(&keys)
	helpers.RenderJSON(w, write, http.StatusOK)
}

// ShareKey relasi antara user dan key
func (h *Handler) ShareKey(w http.ResponseWriter, r *http.Request) {
	//get tribe uint64
	params := mux.Vars(r)
	keyUint, err := strconv.ParseUint(params["key_id"], 10, 32)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][ShareKey][ParseUint]: %s\n", err)
		helpers.SendError(w, "Error when trying to share key", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][ShareKey][ReadBody]: %s\n", err)
		helpers.SendError(w, "Error when trying to share key", http.StatusBadRequest)
		return
	}

	var assign Assign
	//read body, get user id
	if err = json.Unmarshal(body, &assign); err != nil {
		fmt.Printf("[crud_key_handler.go][ShareKey][UnmarshalJSON]: %s\n", err)
		helpers.SendError(w, "Error when trying to share key", http.StatusBadRequest)
		return
	}

	// Get User ID
	uid, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.SendError(w, "error UID extraction", http.StatusInternalServerError)
		return
	}

	var key models.Key
	h.DB.First(&key, uint(keyUint))

	if role < 1 && uint(uid) != key.UserID { // Get user own key
		helpers.SendError(w, "You are not authorized for this request", http.StatusUnauthorized)
		return
	}

	if err = h.DB.Model(&key).Association("Shares").Append(models.KeyShares{UserID: assign.UID, KeyID: uint(keyUint)}).Error; err != nil {
		helpers.SendError(w, "Error when trying to share key", http.StatusBadRequest)
		return
	}

	helpers.SendOK(w, "Shared Key Successfully")
}

// RevokeShare to delete row in sharing key
func (h *Handler) RevokeShare(w http.ResponseWriter, r *http.Request) {
	//get tribe uint64
	params := mux.Vars(r)
	keyUint, err := strconv.ParseUint(params["key_id"], 10, 32)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][RevokeShare][ParseUint]: %s\n", err)
		helpers.SendError(w, "Error when trying to revoke share key", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][RevokeShare][ReadBody]: %s\n", err)
		helpers.SendError(w, "Error when trying to revoke share key", http.StatusBadRequest)
		return
	}

	var assign Assign
	//read body, get user id
	if err = json.Unmarshal(body, &assign); err != nil {
		fmt.Printf("[crud_key_handler.go][RevokeShare][UnmarshalJSON]: %s\n", err)
		helpers.SendError(w, "Error when trying to revoke share key", http.StatusBadRequest)
		return
	}

	// Get User ID
	uid, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.SendError(w, "error UID extraction", http.StatusInternalServerError)
		return
	}

	var key models.Key

	h.DB.First(&key, uint(keyUint))

	if role < 1 && uid != uint64(key.UserID) { // Get user own key
		helpers.SendError(w, "You are not authorized for this request", http.StatusUnauthorized)
		return
	}

	if row := h.DB.Where("user_id = ? AND key_id = ?", assign.UID, keyUint).Delete(models.KeyShares{}).RowsAffected; row == 0 {
		helpers.SendError(w, "Key is not shared with anyone or does not exist", http.StatusBadRequest)
		return
	}

	helpers.SendOK(w, "Revoked Key Access Successfully")
}

// GetKeysHandler returns all keys
func (h *Handler) GetKeysHandler(w http.ResponseWriter, r *http.Request) {
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
	//ADD FILTER, //ADD PAGINATION
	var keys []models.Key
	//IF USER IS SUPERADMIN, GET ALL
	h.DB.Preload("Shares").Preload("Conditions").Where("status = ?", "active").Order("created_at desc").Find(&keys)
	//IF USER IS NORMAL USER, GET ALLOWED (to be updated)

	write, _ := json.Marshal(&keys)
	helpers.RenderJSON(w, write, http.StatusOK)
}

// GetUnregisteredKeys returns all keys
func (h *Handler) GetUnregisteredKeys(w http.ResponseWriter, r *http.Request) {
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
	//ADD FILTER, //ADD PAGINATION
	var keys []models.Key
	//IF USER IS SUPERADMIN, GET ALL
	h.DB.Preload("Shares").Preload("Conditions").Where("status = ?", "unregistered").Order("created_at desc").Find(&keys)
	//IF USER IS NORMAL USER, GET ALLOWED (to be updated)

	write, _ := json.Marshal(&keys)
	helpers.RenderJSON(w, write, http.StatusOK)
}
