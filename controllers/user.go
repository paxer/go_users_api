package controllers

import (
	"go_users_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateUser handles HTTP POST request to create a new User
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Create()
	c.JSON(http.StatusCreated, user)
}
