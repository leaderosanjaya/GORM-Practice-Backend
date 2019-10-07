package key

import "github.com/GORM-practice/app/models"

//Creates
func (h *Handler) CreateKey(key models.Key) error {
	//add keys from the user
	//add keys from the user
	var user models.User
	var tribe models.Tribe
	h.DB.First(&user, key.UserID)
	h.DB.First(&tribe, key.TribeID)

	if dbc := h.DB.Create(&key); dbc.Error != nil {
		return dbc.Error
	}

	h.DB.Model(&user).Association("Keys").Append(key)
	h.DB.Model(&tribe).Association("Keys").Append(key)
	return nil
}

//Deletes Key by providing the given Key ID
func (h *Handler) DeleteKey(targetID uint) error {
	//remove keys from the user
	//remove keys from the tribe, edit tribe key count
	if dbc := h.DB.Where("key_id = ?", targetID).Delete(models.Key{}); dbc.Error != nil {
		return dbc.Error
	}
	return nil
}
