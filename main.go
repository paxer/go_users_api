package main

import (
	"go_users_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.POST("/users", func(c *gin.Context) {
		// TODO pass parameters
		models.NewUser("Darth", "Vader").Create()
		c.String(http.StatusCreated, "TODO: return JSON of a new created user")
	})
	return r
}

func main() {
	db := models.SetupDB()
	defer db.Close()
	r := setupRouter()
	r.Run(":8080")
}
