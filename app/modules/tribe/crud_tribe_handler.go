package tribe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"GORM-practice-backend/app/helpers"
	"GORM-practice-backend/app/models"
	"github.com/gorilla/mux"
)

func (h *Handler) CreateTribeHandler(w http.ResponseWriter, r *http.Request) {
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
	}

	tribe := models.Tribe{}
	if err = json.Unmarshal(body, &tribe); err != nil {
		fmt.Printf("[crud_tribe_handler.go][CreateTribeHandler][UnmarshalJSON]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when creating tribe"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	}

	if err = h.CreateTribe(tribe); err != nil {
		fmt.Printf("[crud_tribe_handler.go][CreateTribeHandler][InsertTribe]: %s\n", err)
		message.Status = "Failed"
		message.Message = "Error when creating tribe"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	}
	helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
	return
}

func (h *Handler) DeleteTribeHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	message := JSONMessage{
		Status:  "Success",
		Message: "Deleted Tribe",
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[crud_tribe_handler.go][DeleteTribeHandler][ReadBody]: %s", err)
		message.Status = "Failed"
		message.Message = "Error while deleting"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	tribeDel := TribeDel{}
	err = json.Unmarshal(body, &tribeDel)
	if err != nil {
		fmt.Printf("[crud_tribe_handler.go][DeleteTribeHandler][UnmarshalJSON]: %s", err)
		message.Status = "Failed"
		message.Message = "Error while deleting"
		status = http.StatusBadRequest
		helpers.RenderJSON(w, helpers.MarshalJSON(message), status)
		return
	}

	if err = h.DeleteTribe(tribeDel.UID); err != nil {
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

func (h *Handler) GetTribeByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var tribe models.Tribe
	h.DB.Preload("Members").Preload("Keys").First(&tribe, params["tribe_id"])
	json.NewEncoder(w).Encode(&tribe)
}
