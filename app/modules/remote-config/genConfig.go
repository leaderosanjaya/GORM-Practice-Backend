package remoteconfig

import (
	"encoding/json"
	"fmt"
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

func (h *Handler) GetConditionData() ([]models.Condition, error) {
	var conditions []models.Condition
	if dbc := h.DB.Find(&conditions); dbc.Error != nil {
		return conditions, dbc.Error
	}
	return conditions, nil
}

// ParseConfig parse the config, return config and error
func (h *Handler) ParseConfig(keys []models.Key, conditions []models.Condition) (Config, error) {
	var config Config
	config.Conditions = []Condition{}
	config.Parameters = map[string]Parameter{}
	for _, condition := range conditions {
		c := Condition{}
		c.Name = condition.ConditionName
		c.Expression = condition.Expression
		c.TagColor = condition.TagColor
		config.Conditions = append(config.Conditions, c)
	}
	for _, key := range keys {
		p := Parameter{DefaultValue: DefaultValue{Value: key.KeyValue}, Description: key.Description}
		var conditionAssign []models.ConditionAssign
		p.ConditionalValues = map[string]Conditionals{}
		h.DB.Where("key_id = ?", key.ID).Find(&conditionAssign)
		for _, v := range conditionAssign {
			condition := models.Condition{}
			h.DB.First(&condition, v.ConditionID)
			p.ConditionalValues[condition.ConditionName] = Conditionals{Value: v.Value}
		}
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
	fmt.Println("[GENCONFIG] Got Key Data")
	conditions, err := h.GetConditionData()
	if err != nil {
		return err
	}
	fmt.Println("[GENCONFIG] Got Condition Data")
	config, err := h.ParseConfig(keys, conditions)
	if err != nil {
		return err
	}
	fmt.Println("[GENCONFIG] Parsed Data")
	err = h.WriteConfig(config)
	if err != nil {
		return err
	}
	return nil
}
