package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Helps Mux to log
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here

		start := time.Now().UTC()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
		end := time.Now().UTC()
		lat := end.Sub(start)

		log.Printf("[API CALL] ROUTE:'%s'[METHOD: %s] %s", r.RequestURI, r.Method, lat)
	})
}

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

// REFACTOR THIS
// JSONMessage struct
type JSONMessage struct {
	Status boolean `json:"status,omitempty"`
	// ErrorCode string  `json:"errorCode,omitempty"`
	Message string `json:"message,omitempty"`
}

func SendOK(w http.ResponseWriter, msg string) {
	message := JSONMessage{
		Status:  true,
		Message: msg,
	}
	RenderJSON(w, MarshalJSON(message), http.StatusOK)
}

func SendError(w http.ResponseWriter, msg string, errCode int) {
	message := JSONMessage{
		Status:  false,
		Message: msg,
	}
	RenderJSON(w, MarshalJSON(message), errCode)
}
