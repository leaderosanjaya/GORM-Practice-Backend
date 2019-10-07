package user

import (
	"github.com/GORM-practice/app/models"

	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) InsertUser(user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	if dbc := h.DB.Create(&user); dbc.Error != nil {
		return dbc.Error
	}
	return nil
}

func (h *Handler) DeleteUser(targetID uint) error {
	if dbc := h.DB.Where("user_id = ?", targetID).Delete(models.User{}); dbc.Error != nil {
		return dbc.Error
	}
	return nil
}
