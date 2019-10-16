package remoteconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/GORM-practice/app/models"
	"github.com/jinzhu/gorm"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

// GetEtag return etag string and error
func (h *Handler) GetEtag() (string, error) {
	//Set up new Client HTTP
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, h.RemoteConfigURL, nil)
	if err != nil {
		return "", err
	}

	// Set Authorization Header
	req.Header.Set("Authorization", "Bearer "+h.Token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	// if resp.Status is 200
	if resp.StatusCode == http.StatusOK {
		return resp.Header["Etag"][0], nil
	}
	return "", fmt.Errorf("Bad Response E-Tag: %d", resp.StatusCode)
}

// GetToken return error
func (h *Handler) GetToken() error {
	var c = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}
	json.Unmarshal([]byte(h.CredentialsFile), &c)
	config := &jwt.Config{
		Email:      c.Email,
		PrivateKey: []byte(c.PrivateKey),
		Scopes: []string{
			"https://www.googleapis.com/auth/firebase.remoteconfig",
		},
		TokenURL: google.JWTTokenURL,
	}
	token, err := config.TokenSource(oauth2.NoContext).Token()
	if err != nil {
		return err
	}
	h.Token = token
	return nil
}

// GetBody return body string and error
func (h *Handler) GetBody() error {
	//Set up new Client HTTP
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, h.RemoteConfigURL, nil)
	if err != nil {
		return err
	}

	// Set Authorization Header
	req.Header.Set("Authorization", "Bearer "+h.Token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	// if resp.Status is 200
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Successfully retrieved latest config")
		// Read response body
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return err
		}

		// Write the response body into config.json
		_ = ioutil.WriteFile("config.json", bodyBytes, 0644)
		return nil
	}
	return fmt.Errorf("Bad Response E-Tag: %d", resp.StatusCode)
}

// PushToDB init db with existing data in firebase
func (h *Handler) PushToDB() error {
	var config Config

	jsonFile, _ := os.Open(os.Getenv("configFile"))
	byteValue, _ := ioutil.ReadAll(jsonFile)

	err := json.Unmarshal(byteValue, &config)
	if err != nil {
		return err
	}

	for k, v := range config.Parameters {
		key := models.Key{}
		if err := h.DB.Where("key_name = ?", k).First(&key).Error; gorm.IsRecordNotFoundError(err) {
			key.KeyName = k
			key.KeyValue = v.DefaultValue.Value
			key.Description = v.Description
			key.Status = "unregistered"

			fmt.Printf("[REMOTE CONFIG INIT] Added key %s to database as unregistered\n", k)

			h.DB.Table("keys").Create(&key)
		}
	}

	return nil
}

// Init to init remote config
func (h *Handler) Init() error {
	h.CredentialsFile = os.Getenv("GOOGLE_CREDENTIALS")
	h.ConfigFile = os.Getenv("configFile")
	h.ProjectID = os.Getenv("PROJECT_ID")
	baseURL := "https://firebaseremoteconfig.googleapis.com"
	remoteConfigEndpoint := fmt.Sprintf("v1/projects/%s/remoteConfig", h.ProjectID)
	h.RemoteConfigURL = fmt.Sprintf("%s/%s", baseURL, remoteConfigEndpoint)

	if err := h.GetToken(); err != nil {
		return err
	}

	if err := h.GetBody(); err != nil {
		return err
	}

	if err := h.PushToDB(); err != nil {
		return err
	}

	fmt.Println("Remote Config Init Successful")
	return nil
}
