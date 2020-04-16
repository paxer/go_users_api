package models

import "time"

// User model, use NewUser to create a new instance
type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
