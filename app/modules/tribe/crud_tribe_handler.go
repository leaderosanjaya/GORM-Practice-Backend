package tribe

import (
	"GORM-practice-backend/app/modules/auth"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"GORM-practice-backend/app/helpers"
	"GORM-practice-backend/app/models"

	"github.com/gorilla/mux"
)

// CreateTribeHandler to handle createtribe
func (h *Handler) CreateTribeHandler(w http.ResponseWriter, r *http.Request) {
	// Get User ID
	_, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"error UID extraction",
		}`), http.StatusInternalServerError)
		return
	}
	if role < 1 {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"Request denied, super admin only",
		}`), http.StatusForbidden)
		return
	}

	status := http.StatusOK
	message := JSONMessage{
		Status:  "Success",
		Message: "Created Tribe",
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_tribe_handler.go][CreateTribeHandler][ReadBody]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when creating tribe"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	tribe := models.Tribe{}
	if err = json.Unmarshal(body, &tribe); err != nil {
		fmt.Printf("[crud_tribe_handler.go][CreateTribeHandler][UnmarshalJSON]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when creating tribe"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	var lead models.User
	if row := h.DB.Where("user_id = ?", tribe.LeadID).First(&lead); row.RowsAffected == 0 {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"Lead does not exist",
		}`), http.StatusBadRequest)
		return
	}

	if err = h.CreateTribe(tribe); err != nil {
		fmt.Printf("[crud_tribe_handler.go][CreateTribeHandler][InsertTribe]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when creating tribe"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}
	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	return
}

//DeleteTribeHandler handle tribe deletion
func (h *Handler) DeleteTribeHandler(w http.ResponseWriter, r *http.Request) {
	// Get User ID
	_, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"error UID extraction",
		}`), http.StatusInternalServerError)
		return
	}
	if role < 1 {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"Request denied, super admin only",
		}`), http.StatusForbidden)
		return
	}

	status := http.StatusOK
	message := JSONMessage{
		Status:  "Success",
		Message: "Deleted Tribe",
	}

	params := mux.Vars(r)

	targetUint, err := strconv.ParseUint(params["tribe_id"], 10, 32)
	if err != nil {
		fmt.Printf("[crud_tribe_handler.go][DeleteTribeHandler][ParseUint]: %s", err)
		message.Status = "Failed"
		message.Message = "Error while deleting"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	if err = h.DeleteTribe(uint(targetUint)); err != nil {
		fmt.Printf("[crud_tribe_handler.go][DeleteTribeHandler][DeleteTribe]: %s", err)
		message.Status = "Failed"
		message.Message = "Error while deleting"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	return
}

// GetTribeByID get tribe by id
func (h *Handler) GetTribeByID(w http.ResponseWriter, r *http.Request) {
	// Get User ID
	_, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"error UID extraction",
		}`), http.StatusInternalServerError)
		return
	}
	if role < 1 {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"Request denied, super admin only",
		}`), http.StatusForbidden)
		return
	}

	params := mux.Vars(r)
	var tribe models.Tribe
	h.DB.Preload("Members").Preload("Keys").First(&tribe, params["tribe_id"])
	json.NewEncoder(w).Encode(&tribe)
}

// AssignUser assign user in tribe by lead
func (h *Handler) AssignUser(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	message := JSONMessage{
		Status:  "Success",
		Message: "Assigned User to Tribe",
	}

	//get tribe uint64
	params := mux.Vars(r)
	tribeUint, err := strconv.ParseUint(params["tribe_id"], 10, 32)
	if err != nil {
		fmt.Printf("[crud_tribe_handler.go][AssignUser][ParseUint]: %s", err)
		message.Status = "Failed"
		message.Message = "Error when trying to add user"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_tribe_handler.go][AssignUser][ReadBody]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when creating tribe"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	var assign Assign
	//read body, get user id
	if err = json.Unmarshal(body, &assign); err != nil {
		fmt.Printf("[crud_tribe_handler.go][AssignUser][UnmarshalJSON]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when creating tribe"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	var user models.User
	if row := h.DB.First(&user, assign.UID); row.RowsAffected == 0 {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"user does not exist",
		}`), http.StatusBadRequest)
		return
	}

	var tribe models.Tribe
	if row := h.DB.First(&tribe, uint(tribeUint)); row.RowsAffected == 0 {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"tribe does not exist",
		}`), http.StatusBadRequest)
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

	if role < 1 && uint64(tribe.LeadID) != uid {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"Request denied, tribe lead or super admin only",
		}`), http.StatusForbidden)
		return
	}

	h.DB.Model(&tribe).Association("Members").Append(models.TribeAssign{UserID: assign.UID, TribeID: uint(tribeUint)})

	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	return
}

// RemoveAssign remove user from tribe by lead
func (h *Handler) RemoveAssign(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	message := JSONMessage{
		Status:  "Success",
		Message: "Removed user from Tribe successfully",
	}

	//get tribe uint64
	params := mux.Vars(r)
	tribeUint, err := strconv.ParseUint(params["tribe_id"], 10, 32)
	if err != nil {
		fmt.Printf("[crud_tribe_handler.go][RemoveAssign][ParseUint]: %s", err)
		message.Status = "Failed"
		message.Message = "Error when trying to remove user"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_tribe_handler.go][RemoveAssign][ReadBody]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when trying to remove user"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	var assign Assign
	//read body, get user id
	if err = json.Unmarshal(body, &assign); err != nil {
		fmt.Printf("[crud_tribe_handler.go][RemoveAssign][UnmarshalJSON]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when trying to remove user"
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

	var tribe models.Tribe
	h.DB.First(&tribe, uint(tribeUint))

	if role < 1 && uid != uint64(tribe.LeadID) {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"Request denied, tribe lead or super admin only",
		}`), http.StatusForbidden)
		return
	}

	if row := h.DB.Where("user_id = ? AND tribe_id = ?", assign.UID, tribeUint).Delete(models.TribeAssign{}); row.RowsAffected == 0 {
		helpers.RenderJSON(w, []byte(`
		{
			"message":"User does not exist in this tribe",
		}`), http.StatusForbidden)
		return
	}

	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	return
}
