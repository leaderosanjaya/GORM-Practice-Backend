package user

import (
	"errors"

	"github.com/GORM-practice/app/models"
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
	if row := h.DB.Where("user_id = ?", targetID).Delete(models.User{}).RowsAffected; row == 0 {
		return errors.New("User does not exist in DB")
	}
	return nil
}
