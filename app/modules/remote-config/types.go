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
	RemoteConfigUrl string
}

type DefaultValue struct {
	Value string `json:"value"`
}

type Parameter struct {
	DefaultValue DefaultValue `json:"defaultValue"`
}

type Config struct {
	Parameters map[string]Parameter `json:"parameters"`
}
