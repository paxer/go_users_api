package models

import (
	"errors"
	"time"
)

// User model, use NewUser to create a new instance
type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new User{} instance
func NewUser(name string, email string) (*User, error) {
	var err error
	if name == "" {
		err = errors.New("name can not be blank")
		return nil, err
	}

	if email == "" {
		err = errors.New("email can not be blank")
		return nil, err
	}

	return &User{Name: name, Email: email}, err
}

// FindUserByID finds User by id and return
func FindUserByID(id int) *User {
	var user User
	DB.First(&user, id)
	return &user
}

// Create creates a new User record in the database
func (u *User) Create() {
	DB.NewRecord(u)
	DB.Create(&u)
}
