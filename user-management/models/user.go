package models

import (
	"log"
	"net/http"
	"user-management/middleware"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	Username   string `json:"username" gorm:"unique"`
	Email      string `json:"email" gorm:"unique"`
	Password   string `json:"password"`
}

func (user *User) GetPassword() string {
	return user.Password
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Println("Error hashing password")
		return err
	}
	user.Password = string(bytes)
	return nil
}
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		log.Println("Error - Passwords don't match")
		return middleware.NewError(http.StatusUnauthorized, "Error - Incorrect Credentials")
	}
	return nil
}
