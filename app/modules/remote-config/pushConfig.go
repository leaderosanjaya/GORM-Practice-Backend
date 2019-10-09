package remoteconfig

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
)

// FUNCTION PUBLISH
// Accepts oauth2 token object & LATEST Etag string
// Doesnt return anything.... should it tho?

// PushData push data, return error
func (h *Handler) PushData(token *oauth2.Token, Etag string) error {
	//Attempt to read config file data
	data, err := ioutil.ReadFile(h.ConfigFile)
	if err != nil {
		return err
		// log.Fatalf("Error retrieving data: %v\n", err)
	}

	//Set up new Client HTTP
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPut, h.RemoteConfigURL, bytes.NewReader(data))
	if err != nil {
		return err
		// log.Fatalf("Error: %v\n", err)
	}
	//To Access Objects
	// fmt.Printf("%+v", token)

	// Set Authorization Header
	req.Header.Set("Authorization", "Bearer "+h.Token.AccessToken)
	req.Header.Add("Content-Type", "application/json; UTF-8")
	req.Header.Add("If-Match", Etag)
	resp, err := client.Do(req)
	if err != nil {
		// log.Fatalf("Error: %v\n", err)
		return err
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Successfully Pushed new config file")
		// fmt.Println("Etag from server: {%s}", resp.Header["Etag"][0])
	} else {
		fmt.Println("Failed to push new config file")
	}
	return nil
}

// PublishConfig publish the config, return error
func (h *Handler) PublishConfig() error {
	//Generates the config
	err := h.GenConfig()
	if err != nil {
		return err
	}
	fmt.Println("Generated Config")
	//GET ETAG
	eTag, err := h.GetEtag()
	if err != nil {
		return err
	}
	fmt.Println("Got e-Tag: ", eTag)
	err = h.PushData(h.Token, eTag)
	if err != nil {
		return err
	}
	return nil
}
