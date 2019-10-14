package tribe

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

// CreateTribeHandler to handle createtribe
func (h *Handler) CreateTribeHandler(w http.ResponseWriter, r *http.Request) {
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_tribe_handler.go][CreateTribeHandler][ReadBody]: %s\n", err)
		helpers.SendError(w, "error creating tribe", http.StatusBadRequest)
		return
	}

	tribe := models.Tribe{}
	if err = json.Unmarshal(body, &tribe); err != nil {
		fmt.Printf("[crud_tribe_handler.go][CreateTribeHandler][UnmarshalJSON]: %s\n", err)
		helpers.SendError(w, "error creating tribe", http.StatusBadRequest)
		return
	}

	var lead models.User
	if row := h.DB.Where("user_id = ?", tribe.LeadID).First(&lead); row.RowsAffected == 0 {
		helpers.SendError(w, "lead does not exist", http.StatusBadRequest)
		return
	}

	if err = h.CreateTribe(tribe); err != nil {
		fmt.Printf("[crud_tribe_handler.go][CreateTribeHandler][InsertTribe]: %s\n", err)
		helpers.SendError(w, "error creating tribe", http.StatusBadRequest)
		return
	}
	helpers.SendOK(w, "tribe created")
	return
}

//DeleteTribeHandler handle tribe deletion
func (h *Handler) DeleteTribeHandler(w http.ResponseWriter, r *http.Request) {
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

	targetUint, err := strconv.ParseUint(params["tribe_id"], 10, 32)
	if err != nil {
		fmt.Printf("[crud_tribe_handler.go][DeleteTribeHandler][ParseUint]: %s", err)
		helpers.SendError(w, "error deleting tribe", http.StatusBadRequest)

		return
	}

	if err = h.DeleteTribe(uint(targetUint)); err != nil {
		fmt.Printf("[crud_tribe_handler.go][DeleteTribeHandler][DeleteTribe]: %s", err)
		helpers.SendError(w, "error deleting tribe", http.StatusBadRequest)
		return
	}

	helpers.SendOK(w, "tribe deleted")
	return
}

// GetTribeByID get tribe by id
func (h *Handler) GetTribeByID(w http.ResponseWriter, r *http.Request) {
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
	var tribe models.Tribe
	h.DB.Preload("Members").Preload("Keys").First(&tribe, params["tribe_id"])
	write, _ := json.Marshal(&tribe)
	helpers.RenderJSON(w, write, http.StatusOK)
}

// AssignUser assign user in tribe by lead
func (h *Handler) AssignUser(w http.ResponseWriter, r *http.Request) {

	//get tribe uint64
	params := mux.Vars(r)
	tribeUint, err := strconv.ParseUint(params["tribe_id"], 10, 32)
	if err != nil {
		fmt.Printf("[crud_tribe_handler.go][AssignUser][ParseUint]: %s", err)
		helpers.SendError(w, "error assign user", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_tribe_handler.go][AssignUser][ReadBody]: %s\n", err)
		helpers.SendError(w, "error assign user", http.StatusBadRequest)
		return
	}

	var assign Assign
	//read body, get user id
	if err = json.Unmarshal(body, &assign); err != nil {
		fmt.Printf("[crud_tribe_handler.go][AssignUser][UnmarshalJSON]: %s\n", err)
		helpers.SendError(w, "error assign user", http.StatusBadRequest)
		return
	}

	var user models.User
	if row := h.DB.First(&user, assign.UID); row.RowsAffected == 0 {
		helpers.SendError(w, "user does not exist", http.StatusBadRequest)
		return
	}

	var tribe models.Tribe
	if row := h.DB.First(&tribe, uint(tribeUint)); row.RowsAffected == 0 {
		helpers.SendError(w, "tribe does not exist", http.StatusBadRequest)
		return
	}

	// Get User ID
	uid, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.SendError(w, "error uid extraction", http.StatusInternalServerError)
		return
	}

	if role < 1 && uint64(tribe.LeadID) != uid {
		helpers.SendError(w, "super admin access only", http.StatusForbidden)
		return
	}

	h.DB.Model(&tribe).Association("Members").Append(models.TribeAssign{UserID: assign.UID, TribeID: uint(tribeUint), Platform: assign.PlatformID})

	helpers.SendOK(w, "user assigned")
	return
}

// RemoveAssign remove user from tribe by lead
func (h *Handler) RemoveAssign(w http.ResponseWriter, r *http.Request) {
	
	//get tribe uint64
	params := mux.Vars(r)
	tribeUint, err := strconv.ParseUint(params["tribe_id"], 10, 32)
	if err != nil {
		fmt.Printf("[crud_tribe_handler.go][RemoveAssign][ParseUint]: %s", err)
		helpers.SendError(w, "error remove user", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_tribe_handler.go][RemoveAssign][ReadBody]: %s\n", err)
		helpers.SendError(w, "error remove user", http.StatusBadRequest)
		return
	}

	var assign Assign
	//read body, get user id
	if err = json.Unmarshal(body, &assign); err != nil {
		fmt.Printf("[crud_tribe_handler.go][RemoveAssign][UnmarshalJSON]: %s\n", err)
		helpers.SendError(w, "error remove user", http.StatusBadRequest)
		return
	}

	// Get User ID
	uid, role, err := auth.ExtractTokenUID(r)
	if err != nil {
		helpers.SendError(w, "error uid extraction", http.StatusInternalServerError)
		return
	}

	var tribe models.Tribe
	h.DB.First(&tribe, uint(tribeUint))

	if role < 1 && uid != uint64(tribe.LeadID) {
		helpers.SendError(w, "tribe lead or super admin access only", http.StatusForbidden)
		return
	}

	if row := h.DB.Where("user_id = ? AND tribe_id = ?", assign.UID, tribeUint).Delete(models.TribeAssign{}); row.RowsAffected == 0 {
		helpers.SendError(w, "user does not exist", http.StatusBadRequest)
		return
	}

	helpers.SendOK(w, "removed user")
	return
}
