package key

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
		return
	}

	key := models.Key{}
	if err = json.Unmarshal(body, &key); err != nil {
		fmt.Printf("[crud_key_handler.go][CreateKeyHandler][UnmarshalJSON]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when creating key"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	// Get user ID from token
	uid, _, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"error UID extraction",
		}`), http.StatusInternalServerError)
		return
	}

	key.UserID = uint(uid)

	var user models.TribeAssign
	h.DB.Where("user_id = ?", uint(uid)).First(&user)
	key.TribeID = user.TribeID

	if err = h.CreateKey(key); err != nil {
		fmt.Printf("[crud_key_handler.go][CreateKeyHandler][InsertTribe]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when creating key"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	h.PushRemoteConfig()
	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
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
		fmt.Printf("[crud_key_handler.go][DeleteKeyHandler][ParseUint]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error while deleting"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	// Get user ID from token
	uid, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"error UID extraction",
		}`), http.StatusInternalServerError)
		return
	}

	var key models.Key
	row := h.DB.First(&key, params["key_id"]).RowsAffected
	if row == 0 {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"Key does not exist in DB",
		}`), http.StatusForbidden)
		return
	}

	if role < 1 && uint64(key.UserID) != uid {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"Failed to delete, you are not the owner of this key",
		}`), http.StatusForbidden)
		return
	}

	if err = h.DeleteKey(uint(targetUint)); err != nil {
		fmt.Printf("[crud_key_handler.go][DeleteKeyHandler][DeleteTribe]: %s\n", err)
		message.Status = "Failed"
		message.Message = err.Error()
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	h.PushRemoteConfig()

	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
}

//GetKeyByID by user
func (h *Handler) GetKeyByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var key models.Key

	// Get user ID from token
	uid, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"error UID extraction",
		}`), http.StatusInternalServerError)
		return
	}

	row := h.DB.Preload("Shares").First(&key, params["key_id"]).RowsAffected
	if row == 0 {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"Key does not exist",
		}`), http.StatusBadRequest)
		return
	}

	if key.UserID != uint(uid) && role < 1 {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"you are not the owner of this key",
		}`), http.StatusUnauthorized)
		return
	}

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

	uid, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"error UID extraction",
		}`), http.StatusInternalServerError)
		return
	}

	if role < 1 && uint64(key.UserID) != uid { // Get user own key
		helpers.RenderJSON(w, []byte(`
		{
			"message":"you are not authorized to request",
		}`), http.StatusUnauthorized)
		return
	}

	//read edit info
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][UpdateKeyByID][ReadBody]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when updating key"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	updateKey := models.Key{}
	if err = json.Unmarshal(body, &updateKey); err != nil {
		fmt.Printf("[crud_key_handler.go][UpdateKeyByID][UnmarshalJSON]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when updating key"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	updateValue(&updateKey, &key)
	h.DB.Save(&key)

	h.PushRemoteConfig()
	message = JSONMessage{
		Status:  "Success",
		Message: "Updated Key",
	}
	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
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

	if paramUserID, _ := strconv.ParseUint(fmt.Sprintf("%s", params["user_id"]), 10, 32); paramUserID != uid && role < 1 { // Get user own key
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

	uid, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"error UID extraction",
		}`), http.StatusInternalServerError)
		return
	}

	paramTribeID, _ := strconv.ParseUint(fmt.Sprintf("%s", params["tribe_id"]), 10, 32)

	if role < 1 {
		var userTribes []models.TribeAssign
		if row := h.DB.Where("user_id = ?", uint(uid)).Find(&userTribes); row.RowsAffected == 0 {
			helpers.RenderJSON(w, []byte(`
			{
				"message":"user does not exist",
			}`), http.StatusBadRequest)
			return
		}
		if len(userTribes) == 0 {
			helpers.RenderJSON(w, []byte(`
			{
				"message":"user is not from this tribe, request denied",
			}`), http.StatusForbidden)
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

	h.DB.Preload("Shares").Where("tribe_id = ?", params["tribe_id"]).Find(&keys)
	json.NewEncoder(w).Encode(&keys)
}

// ShareKey relasi antara user dan key
func (h *Handler) ShareKey(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	message := JSONMessage{
		Status:  "Success",
		Message: "Shared Key Successfully",
	}

	//get tribe uint64
	params := mux.Vars(r)
	keyUint, err := strconv.ParseUint(params["key_id"], 10, 32)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][ShareKey][ParseUint]: %s\n", err)
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
		return
	}

	var assign Assign
	//read body, get user id
	if err = json.Unmarshal(body, &assign); err != nil {
		fmt.Printf("[crud_key_handler.go][ShareKey][UnmarshalJSON]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when trying to share key"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	// Get User ID
	uid, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"error UID extraction",
		}`), http.StatusInternalServerError)
		return
	}

	var key models.Key
	h.DB.First(&key, uint(keyUint))

	if role < 1 && uint(uid) != key.UserID { // Get user own key
		helpers.RenderJSON(w, []byte(`
		{
			"message":"you are not authorized to request",
		}`), http.StatusUnauthorized)
		return
	}

	if err = h.DB.Model(&key).Association("Shares").Append(models.KeyShares{UserID: assign.UID, KeyID: uint(keyUint)}).Error; err != nil {
		message.Status = "Failed"
		message.Message = fmt.Sprintf("Error when trying to share key, check your request again \n %s", err.Error())
		helpers.RenderJSON(w, helpers.MarshalJSON(message), http.StatusUnauthorized)
		return
	}

	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
}

// RevokeShare to delete row in sharing key
func (h *Handler) RevokeShare(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	message := JSONMessage{
		Status:  "Success",
		Message: "Revoked Key Share Access Successfully",
	}

	//get tribe uint64
	params := mux.Vars(r)
	keyUint, err := strconv.ParseUint(params["key_id"], 10, 32)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][RevokeShare][ParseUint]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when trying to revoke share key"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][RevokeShare][ReadBody]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when trying to revoke share key"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	var assign Assign
	//read body, get user id
	if err = json.Unmarshal(body, &assign); err != nil {
		fmt.Printf("[crud_key_handler.go][RevokeShare][UnmarshalJSON]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when trying to revoke share key"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	// Get User ID
	uid, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"error UID extraction",
		}`), http.StatusInternalServerError)
		return
	}

	var key models.Key

	h.DB.First(&key, uint(keyUint))

	if role < 1 && uid != uint64(key.UserID) { // Get user own key
		helpers.RenderJSON(w, []byte(`
		{
			"message":"you are not authorized to request",
		}`), http.StatusUnauthorized)
		return
	}

	if row := h.DB.Where("user_id = ? AND key_id = ?", assign.UID, keyUint).Delete(models.KeyShares{}).RowsAffected; row == 0 {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"Key is not shared with anyone, or key does not exist",
		}`), http.StatusBadRequest)
		return
	}

	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
}
