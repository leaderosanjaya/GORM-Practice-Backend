package key

import (
	"github.com/GORM-practice/app/models"
)

//CreateKey create key
func (h *Handler) CreateKey(key models.Key) error {
	//Initialize User and Tribe variable
	var user models.User
	var tribe models.Tribe
	//Get related user & tribe
	h.DB.First(&user, key.UserID)
	h.DB.First(&tribe, key.TribeID)

	//Execute Create key
	if dbc := h.DB.Create(&key); dbc.Error != nil {
		return dbc.Error
	}

	//Associate new key to related user and tribe
	h.DB.Model(&user).Association("Keys").Append(key)
	h.DB.Model(&tribe).Association("Keys").Append(key)
	return nil
}

//DeleteKey by providing the given Key ID
func (h *Handler) DeleteKey(targetID uint) error {
	//Get target and execute delete
	if err := h.DB.Where("key_id = ?", targetID).Delete(models.Key{}).Error; err != nil {
		return err
	}
	return nil
}

func updateValue(updateKey *models.Key, key *models.Key) {
	//Function for update key
	if updateKey.KeyName != "" {
		key.KeyName = updateKey.KeyName
	}
	if updateKey.KeyValue != "" {
		key.KeyValue = updateKey.KeyValue
	}
	if updateKey.KeyType != "" {
		key.KeyType = updateKey.KeyType
	}
	if updateKey.Description != "" {
		key.Description = updateKey.Description
	}
	if updateKey.Platform != "" {
		key.Platform = updateKey.Platform
	}
	if updateKey.ExpireDate.IsZero() {
		key.ExpireDate = updateKey.ExpireDate
	}
	if updateKey.UserID != 0 {
		key.UserID = updateKey.UserID
	}
	if updateKey.TribeID != 0 {
		key.TribeID = updateKey.TribeID
	}
	if updateKey.AppVersion != "" {
		key.AppVersion = updateKey.AppVersion
	}
	if updateKey.Status != "" {
		key.Status = updateKey.Status
	}
	return
}
