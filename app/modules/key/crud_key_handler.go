package key

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"rc-practice-backend/app/helpers"

	"GORM-practice-backend/app/models"
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][DeleteKeyHandler][ReadBody]: %s", err)
		message.Status = "Failed"
		message.Message = "Error while deleting"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	keyDel := Del{}
	err = json.Unmarshal(body, &keyDel)
	if err != nil {
		fmt.Printf("[crud_key_handler.go][DeleteKeyHandler][UnmarshalJSON]: %s", err)
		message.Status = "Failed"
		message.Message = "Error while deleting"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	if err = h.DeleteKey(keyDel.UID); err != nil {
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

//Getkey by user
