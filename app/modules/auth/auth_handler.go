package auth

import (
	"log"
	"time"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/GORM-practice/app/helpers"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type key string

const user key = "user"

// JwtVerify Verify jwt token for every request
func JwtVerify(next http.Handler) http.Handler {
	return (http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var header = r.Header.Get("Cookie")
		if header == "" {
			helpers.SendError(w, "Error: Found no token in header", http.StatusBadRequest)
			return
		}

		headerSplit := strings.Split(header, "=")
		if len(headerSplit) != 2 {
			helpers.SendError(w, "Missing auth token", http.StatusBadRequest)
			return
		}

		header = headerSplit[1]

		tk := &Token{}

		err := godotenv.Load()
		if err != nil {
			fmt.Printf("[DB Load Env] %s\n", err)
			return
		}

		refreshToken(w, r)

		ctx := context.WithValue(r.Context(), user, tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	}))
}

// Login user login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	cred := Credential{}
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		helpers.SendError(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	if lenPass := len(cred.Password); lenPass < 6 {
		helpers.SendError(w, "Invalid password length, password lenght should be more than 6 character", http.StatusUnprocessableEntity)
		return
	}

	resp := h.FindOne(cred.Email, cred.Password)

	tokenString := resp["token"].(string)
	tokenExpiresAt := time.Unix(resp["expiresAt"].(int64), 0)

	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: tokenString,
		Expires: tokenExpiresAt,
		Path: "/",
	})

	message, _ := json.Marshal(resp)
	helpers.RenderJSON(w, message, http.StatusOK)
	// json.NewEncoder(w).Encode(resp)
}

// ExtractToken to extract token from http request header
func ExtractToken(r *http.Request) string {
	if cookieToken := r.Header.Get("Cookie"); cookieToken!="" {
		return strings.Split(cookieToken, "=")[1]
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
		role, err := strconv.ParseInt(fmt.Sprintf("%.0f", claims["Role"]), 10, 0)
		if err != nil {
			return 0, 0, err
		}
		return uid, role, nil
	}
	return 0, 0, err
}

// ValidateToken to validate token request
func (h *Handler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	tokenString := r.Header.Get("Authorization")
	splitToken := strings.Split(tokenString, " ")
	tokenString = splitToken[1]

	// Initialize a new instance of `Claims`
	claims := Token{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("API_SECRET_KEY")), nil
	})

	if err != nil {
		helpers.SendError(w, "invalid token format", http.StatusUnauthorized)
		return
	}
	if err == jwt.ErrSignatureInvalid {
		helpers.SendError(w, "token signature invalid", http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		helpers.SendError(w, "invalid token, this request has no authorization token", http.StatusUnauthorized)
		return
	}

	helpers.SendOK(w, "token is valid")
}

// refreshToken to refresh token each time jwt Verify
func refreshToken(w http.ResponseWriter, r *http.Request) {
	oldToken := ExtractToken(r)
	if oldToken == "" {
		helpers.SendError(w, "request has no authorization token", http.StatusUnauthorized)
		return
	}

	claims := Token{}
	token, err := jwt.ParseWithClaims(oldToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			helpers.SendError(w, "invalid signature token", http.StatusUnauthorized)
			return
		}
		log.Println("error:::::", err)
		helpers.SendError(w, "failed parsing claims", http.StatusBadRequest)
		return
	}

	if !token.Valid{
		helpers.SendError(w, "token invalid", http.StatusUnauthorized)
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) <= 3*time.Minute {
		expiresAt := time.Now().Add(10 * time.Minute)
		claims.ExpiresAt = expiresAt.Unix()
		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)	
		if err != nil {
			helpers.SendError(w, "failed to refresh token", http.StatusInternalServerError)
			return
		}

		newTokenString, err := newToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
		if err != nil {
			helpers.SendError(w, "failed to refresh token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name: "token",
			Value: newTokenString,
			Expires: expiresAt,
			Path: "/",
		})
		log.Println("Token refreshed")
	}
	return
}

// Logout is to set token expiration time to 
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	oldToken := ExtractToken(r)
	if oldToken == "" {
		helpers.SendError(w, "request has no authorization token", http.StatusUnauthorized)
		return
	}

	claims := Token{}
	token, err := jwt.ParseWithClaims(oldToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			helpers.SendError(w, "invalid signature token", http.StatusUnauthorized)
			return
		}
		log.Println("error:::::", err)
		helpers.SendError(w, "failed parsing claims", http.StatusBadRequest)
		return
	}

	if !token.Valid{
		helpers.SendError(w, "token invalid", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: "deleted",
		Expires: time.Now().Add(-100 * time.Hour),
		Path: "/",
	})
	log.Println("Token expired")

	helpers.SendOK(w, "Logged Out")
	return
}