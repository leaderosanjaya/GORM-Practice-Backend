package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GORM-practice/app/models"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
)

// FindOne verify user email and password
func (h *Handler) FindOne(email, password string) map[string]interface{} {
	user := models.User{}

	if err := h.DB.Where("email= ?", email).First(&user).Error; err != nil {
		resp := map[string]interface{}{"status": false, "message": "email not found"}
		return resp
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println(err)
		resp := map[string]interface{}{"status": false, "message": "credential false"}
		return resp
	}

	// JWT
	expiresAt := time.Now().Add(time.Minute * 10).Unix()

	tk := &Token{
		UserID: user.ID,
		Name:   user.LastName,
		Email:  user.Email,
		Role:   user.Role,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Println(err)
		return nil
	}

	resp := map[string]interface{}{"status": true, "message": "logged in"}
	resp["token"] = tokenString
	resp["user"] = user
	return resp
}
