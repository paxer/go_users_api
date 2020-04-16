package models

import "github.com/jinzhu/gorm"

// User model, use NewUser to create a new instance
type User struct {
	gorm.Model
	Name  string
	Email string
}

// NewUser creates a new User{} instance
func NewUser(name string, email string) *User {
	return &User{Name: name, Email: email}
}

// Create creates a new User record in the database
func (u *User) Create() {
	DB.NewRecord(u)
	DB.Create(&u)
}
