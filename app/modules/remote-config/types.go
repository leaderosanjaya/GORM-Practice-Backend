package remoteconfig

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
)

// Handler handler struct
type Handler struct {
	DB              *gorm.DB
	Token           *oauth2.Token
	CredentialsFile string
	ConfigFile      string
	ProjectID       string
	RemoteConfigURL string
}

// DefaultValue default value json
type DefaultValue struct {
	Value string `json:"value"`
}

// Parameter default value defaultValue
type Parameter struct {
	DefaultValue      DefaultValue            `json:"defaultValue"`
	Description       string                  `json:"description"`
	ConditionalValues map[string]Conditionals `json:"conditionalValues,omitempty"`
}

type Conditionals struct {
	Value string `json:"value"`
}

// Condition struct
type Condition struct {
	Name       string `json:"name"`
	Expression string `json:"expression"`
	TagColor   string `json:"tagColor"`
}

// Config parameters
type Config struct {
	Conditions []Condition          `json:"conditions"`
	Parameters map[string]Parameter `json:"parameters"`
}
