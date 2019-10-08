package auth

import (
	"GORM-practice-backend/app/helpers"
	"GORM-practice-backend/app/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type key string

const user key = "user"

// JwtVerify Verify jwt token for every request
func JwtVerify(next http.Handler) http.Handler {
	return (http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var header = r.Header.Get("x-access-token")

		header = strings.TrimSpace(header)

		if header == "" {
			w.WriteHeader(http.StatusForbidden)
			helpers.RenderJSON(w, []byte(`
			{
				message: "missing auth token",
			}
			`), http.StatusBadRequest)
			return
		}

		tk := &Token{}

		err := godotenv.Load()
		if err != nil {
			fmt.Printf("[DB Load Env] %s\n", err)
			return
		}

		_, err = jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			helpers.RenderJSON(w, []byte(`
			{
				message: "error",
			}
			`), http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), user, tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	}))
}

// Login user login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		resp := map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
	}

	resp := h.FindOne(user.Email, user.Password)
	json.NewEncoder(w).Encode(resp)
}
