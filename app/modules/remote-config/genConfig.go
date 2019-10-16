package remoteconfig

import (
	"encoding/json"
	"io/ioutil"

	"github.com/GORM-practice/app/models"
)

// GetKeyData get key data, return slice of key and error
func (h *Handler) GetKeyData() ([]models.Key, error) {
	var keys []models.Key
	if dbc := h.DB.Where("status != ?", "inactive").Find(&keys); dbc.Error != nil {
		return keys, dbc.Error
	}
	return keys, nil
}

// ParseConfig parse the config, return config and error
func (h *Handler) ParseConfig(keys []models.Key) (Config, error) {
	var config Config
	config.Parameters = map[string]Parameter{}
	for _, key := range keys {
		p := Parameter{DefaultValue: DefaultValue{Value: key.KeyValue}, Description: key.Description}
		config.Parameters[key.KeyName] = p
	}
	return config, nil
}

// WriteConfig write config, return error
func (h *Handler) WriteConfig(config Config) error {
	file, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config.json", file, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GenConfig generate config, return error
func (h *Handler) GenConfig() error {
	keys, err := h.GetKeyData()
	if err != nil {
		return err
	}
	config, err := h.ParseConfig(keys)
	if err != nil {
		return err
	}
	err = h.WriteConfig(config)
	if err != nil {
		return err
	}
	return nil
}
