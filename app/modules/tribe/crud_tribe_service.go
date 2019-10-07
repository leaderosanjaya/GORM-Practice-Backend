package tribe

import "GORM-practice-backend/app/models"

func (h *Handler) CreateTribe(tribe models.Tribe) error {
	//Get tribe lead id
	//Insert tribe to lead user
	var lead models.User
	h.DB.First(&lead, tribe.LeadID)

	if dbc := h.DB.Create(&tribe); dbc.Error != nil {
		return dbc.Error
	}
	h.DB.Model(&lead).Association("Tribes").Append(models.TribeAssign{UserID: lead.ID, TribeID: tribe.ID})
	return nil
}

func (h *Handler) DeleteTribe(targetID uint) error {
	if dbc := h.DB.Where("tribe_id = ?", targetID).Delete(models.Tribe{}); dbc.Error != nil {
		return dbc.Error
	}
	return nil
}
