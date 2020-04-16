package controllers

import (
	"go_users_api/models"
	"net/http"
	"strconv"

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

// ShowUser handles HTTP GET request to find a User by id param
func ShowUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.FindUserByID(id)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser handles HTTP PATCH request to find and update User details
func UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.FindUserByID(id)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Update()
	c.JSON(http.StatusOK, user)
}

// DeleteUser handles HTTP DELETE request to find and delete User record
func DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.FindUserByID(id)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	user.Delete()
	c.JSON(http.StatusOK, gin.H{})
}
