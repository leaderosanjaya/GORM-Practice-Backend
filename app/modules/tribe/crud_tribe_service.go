package tribe

import (
	"errors"

	"github.com/GORM-practice/app/models"
)

// CreateTribe create tribe
func (h *Handler) CreateTribe(tribe tCreate) (uint, error) {
	// Migrate local tribe to gorm struct
	newTribe := models.Tribe{}

	newTribe.TribeName = tribe.TribeName
	newTribe.Description = tribe.Description

	//Get tribe lead id
	//Insert tribe to lead user

	if tribe.LeadID != 0 {
		newTribe.TotalMember = 1
	}

	if dbc := h.DB.Create(&newTribe); dbc.Error != nil {
		return 0, dbc.Error
	}

	if tribe.LeadID != 0 {
		lead := models.User{}
		if err := h.DB.First(&lead, tribe.LeadID); err.RowsAffected == 0 {
			return 0, errors.New("user lead does not exist, tribe created without lead")
		}
		h.DB.Model(&newTribe).Association("Leads").Append(models.TribeLeadAssign{LeadID: lead.ID, TribeID: newTribe.ID})
		h.DB.Model(&lead).Association("Tribes").Append(models.TribeAssign{UserID: lead.ID, TribeID: newTribe.ID})
	}
	return newTribe.ID, nil
}

// DeleteTribe delete tribe
func (h *Handler) DeleteTribe(targetID uint) error {
	if row := h.DB.Where("tribe_id = ?", targetID).Delete(models.Tribe{}); row.RowsAffected == 0 {
		return errors.New("tribe does not exist")
	}
	return nil
}

// UpdateValue update tribe value
func UpdateValue(updateTribe *models.Tribe, tribe *models.Tribe) {
	if updateTribe.TribeName != "" {
		tribe.TribeName = updateTribe.TribeName
	}
	if updateTribe.Description != "" {
		tribe.Description = updateTribe.Description
	}
}

// // GetTribeID returns tribe ID
// func (h *Handler) GetTribeID(tribeName string) (uint, error) {
// 	tribe := models.Tribe{}
// 	if row := h.DB.Where("tribe_name = ?", tribeName).First(&tribe); row.RowsAffected == 0 {
// 		return 0, errors.New("tribe does not exist")
// 	}
// 	return tribe.ID, nil
// }