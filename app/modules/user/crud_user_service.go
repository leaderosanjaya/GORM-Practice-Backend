package user

import (
	"GORM-practice-backend/app/models"

	"golang.org/x/crypto/bcrypt"
)

// InsertUser insert user to db
func (h *Handler) InsertUser(user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	if dbc := h.DB.Create(&user); dbc.Error != nil {
		return dbc.Error
	}
	return nil
}

// DeleteUser delete user from db
func (h *Handler) DeleteUser(targetID uint) error {
	if dbc := h.DB.Where("user_id = ?", targetID).Delete(models.User{}); dbc.Error != nil {
		return dbc.Error
	}
	return nil
}
