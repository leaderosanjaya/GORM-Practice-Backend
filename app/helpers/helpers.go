package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// RenderJSON returns message in JSON to http body
func RenderJSON(w http.ResponseWriter, data []byte, status int) {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// set HTTP respon type to JSON
	w.Header().Set("Content-Type", "application/json")

	// HTTP status (200 OK, 404 Not Found, 500 Internal Server Error, etc.)
	w.WriteHeader(status)

	// The actual data
	w.Write(data)
}

// HashPassword to hash user password
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// ComparePassword to compare user input password with hashed password stored in db
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// MarshalJSON to marshal j
func MarshalJSON(message interface{}) []byte {
	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err.Error())
		jsonData = []byte(err.Error())
	}
	return jsonData
}
