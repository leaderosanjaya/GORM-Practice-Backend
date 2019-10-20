package condition

import (
	"strconv"
	"github.com/gorilla/mux"
	"encoding/json"
	"github.com/GORM-practice/app/models"
	"io/ioutil"
	"github.com/GORM-practice/app/helpers"
	"github.com/labstack/gommon/log"
	"net/http"
)

// CreateConditionHandler to create condition
func (h *Handler) CreateConditionHandler(w http.ResponseWriter, r *http.Request){
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[Create Condition Read Body][Condition]: %s\n", err)
		helpers.SendError(w, "error create condition", http.StatusBadRequest)
		return
	}

	condition := models.Condition{}
	if err := json.Unmarshal(body, &condition); err != nil {
		log.Printf("[Create Condition Unmarshal Body][Condition]: %s\n", err)
		helpers.SendError(w, "error create condition", http.StatusBadRequest)
		return
	}

	if err := h.CreateCondition(condition); err != nil {
		log.Printf("[Create Condition][Condition]: %s\n", err)
		helpers.SendError(w, "error create condition", http.StatusBadRequest)
		return
	}
	
	helpers.SendOK(w, "Condition created successfully")
}

// RetrieveConditionHandler to get condition by condition_id
func (h *Handler) RetrieveConditionHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	
	conditionID, err := strconv.ParseUint(params["condition_id"], 10, 0)
	if err != nil {
		log.Printf("[Retrieve Condition strconv][Condition]: %s\n", err)
		helpers.SendError(w, "failed retrieve condition", http.StatusInternalServerError)
		return
	}

	condition, err := h.RetrieveCondition(uint(conditionID))
	if err != nil {
		log.Printf("[Retrieve Condition][Condition]: %s\n", err)
		helpers.SendError(w, "failed retrieve condition", http.StatusInternalServerError)
		return
	}

	write, _ := json.Marshal(&condition)
	helpers.RenderJSON(w, write, http.StatusOK)
}

// RetrieveConditionsHandler to get all conditions
func (h *Handler) RetrieveConditionsHandler(w http.ResponseWriter, r *http.Request){
	
	conditions, err := h.RetrieveConditions()
	if err != nil {
		log.Printf("[Retrieve Conditions][Condition]: %s\n", err)
		helpers.SendError(w, "failed retrieve condition", http.StatusInternalServerError)
		return
	}

	write, _ := json.Marshal(&conditions)
	helpers.RenderJSON(w, write, http.StatusOK)
}

// UpdateConditionHandler to update condition by condition_id
func (h *Handler) UpdateConditionHandler(w http.ResponseWriter, r *http.Request){
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[update Condition Read Body][Condition]: %s\n", err)
		helpers.SendError(w, "error update condition", http.StatusBadRequest)
		return
	}

	updateCondition := models.Condition{}
	if err := json.Unmarshal(body, &updateCondition); err != nil {
		log.Printf("[update Condition Unmarshal Body][Condition]: %s\n", err)
		helpers.SendError(w, "error update condition", http.StatusBadRequest)
		return
	}

	if err := h.UpdateCondition(updateCondition); err != nil {
		log.Printf("[update Condition][Condition]: %s\n", err)
		helpers.SendError(w, "error update condition", http.StatusInternalServerError)
		return
	}
	
	helpers.SendOK(w, "Condition updated successfully")
}

// DeleteConditionHandler to delete condition by condition_id
func (h *Handler) DeleteConditionHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	conditionID, err := strconv.ParseUint(params["condition_id"], 10, 0)
	if err != nil {
		log.Printf("[Delete Condition strconv][Condition]: %s\n", err)
		helpers.SendError(w, "failed Delete condition", http.StatusInternalServerError)
		return
	}

	if err := h.DeleteCondition(uint(conditionID)); err != nil {
		log.Printf("[Delete Condition][Condition]: %s\n", err)
		helpers.SendError(w, "error Delete condition", http.StatusInternalServerError)
		return
	}
}