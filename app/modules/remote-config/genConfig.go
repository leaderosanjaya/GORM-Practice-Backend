package remoteconfig

import "github.com/GORM-practice/app/models"

//Get all keys
//Get key name, Key Value

func (h *Handler) GetKeyData() ([]models.Key, error) {
	var keys []models.Key
	if dbc := h.DB.Where("status = ?", "active").Find(&keys); dbc.Error != nil {
		return keys, dbc.Error
	}
	return keys, nil
}

func (h *Handler) ParseConfig() error {

}
