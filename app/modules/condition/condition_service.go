package condition

import (
	"errors"
	"github.com/GORM-practice/app/models"
)


// CreateCondition query DB to insert condition to DB
func (h *Handler) CreateCondition(condition models.Condition) error {
	if dbc := h.DB.Create(&condition); dbc.Error != nil {
		return dbc.Error
	}
	return nil
}

// DeleteCondition delete condition from DB
func (h *Handler) DeleteCondition(conditionID uint) error {
	if row := h.DB.Where("condition_id = ?", conditionID).Delete(models.Condition{}).RowsAffected; row == 0 {
		return errors.New("Condition does not exist in DB")
	}
	return nil
}

// UpdateCondition update condition row in DB
func (h *Handler) UpdateCondition(updateCondition models.Condition) error {

	var condition models.Condition
	if row := h.DB.First(&condition, updateCondition.ID).RowsAffected; row == 0 {
		return errors.New("condition does not exist")
	}

	UpdateValue(&updateCondition, &condition)
	if err := h.DB.Save(&condition).Error; err != nil {
		return err
	}
	return nil
}

// RetrieveCondition retrieve condition by condition_id from DB
func (h *Handler) RetrieveCondition(conditionID uint) (models.Condition, error) {
	var condition models.Condition
	if err := h.DB.Where("condition_id = ?", conditionID).First(&condition).Error; err != nil {
		return models.Condition{}, err
	}
	return condition, nil
}

// RetrieveConditions retrieve all conditions from DB
func (h *Handler) RetrieveConditions() ([]models.Condition, error) {
	var conditions []models.Condition
	if err := h.DB.Preload("AffectingKeys").Find(&conditions).Error; err != nil {
		return nil, err
	}
	return conditions, nil
}

// UpdateValue update the updated value in condition update
func UpdateValue(updateCondition, condition *models.Condition) {
	if updateCondition.ConditionName != "" {
		condition.ConditionName = updateCondition.ConditionName
	}
	if updateCondition.Expression != "" {
		condition.Expression = updateCondition.Expression
	}
	if updateCondition.TagColor != "" {
		condition.TagColor = updateCondition.TagColor
	}

	// IMPROVE: add password change for user
}
