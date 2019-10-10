package tribe

import (
	"errors"

	"GORM-practice-backend/app/models"
)

// CreateTribe create tribe
func (h *Handler) CreateTribe(tribe models.Tribe) error {
	//Get tribe lead id
	//Insert tribe to lead user
	var lead models.User

	if err := h.DB.First(&lead, tribe.LeadID); err.RowsAffected == 0 {
		return errors.New("calon lead does not exist")
	}

	if dbc := h.DB.Create(&tribe); dbc.Error != nil {
		return dbc.Error
	}
	h.DB.Model(&lead).Association("Tribes").Append(models.TribeAssign{UserID: lead.ID, TribeID: tribe.ID})
	return nil
}

// DeleteTribe delete tribe
func (h *Handler) DeleteTribe(targetID uint) error {
	if row := h.DB.Where("tribe_id = ?", targetID).Delete(models.Tribe{}); row.RowsAffected == 0 {
		return errors.New("tribe does not exist")
	}
	return nil
}
