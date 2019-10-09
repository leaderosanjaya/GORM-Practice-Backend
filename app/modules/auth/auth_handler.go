package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/GORM-practice/app/helpers"
	"github.com/GORM-practice/app/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type key string

const user key = "user"

// JwtVerify Verify jwt token for every request
func JwtVerify(next http.Handler) http.Handler {
	return (http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var header = r.Header.Get("Authorization")

		if header == "" {
			w.WriteHeader(http.StatusForbidden)
			helpers.RenderJSON(w, []byte(`
			{
				message: "missing auth token",
			}
			`), http.StatusBadRequest)
			return
		}

		headerSplit := strings.Split(header, " ")
		if len(headerSplit) != 2 {
			w.WriteHeader(http.StatusForbidden)
			helpers.RenderJSON(w, []byte(`
			{
				message: "missing auth token",
			}
			`), http.StatusBadRequest)
			return
		}

		header = headerSplit[1]

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
			log.Println(err)
			helpers.RenderJSON(w, []byte(`
			{
				message: "error, no auth token found, or your auth token is false",
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

// ExtractToken to extract token from http request header
func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}

	bearerToken := r.Header.Get("Authorization")
	if token = strings.Split(bearerToken, " ")[1]; token != "" {
		return token
	}

	return ""
}

// ExtractTokenUID extract token from request
func ExtractTokenUID(r *http.Request) (uint64, int64, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing error")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return 0, 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["UserID"]), 10, 32)
		role, err := strconv.ParseInt(fmt.Sprintf("%.0f", claims["Role"]), 10, 32)
		if err != nil {
			return 0, 0, err
		}
		return uid, role, nil
	}
	return 0, 0, err
}
