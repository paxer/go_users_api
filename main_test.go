package main

import (
	"bytes"
	"encoding/json"
	"go_users_api/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
)

func setupDB() *gorm.DB {
	os.Setenv("DB_CONN", "host=localhost port=5432 user=postgres dbname=go_users_api_test password=postgres sslmode=disable")
	db := models.SetupDB()
	return db
}

func Test_PingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected code 200 got %d", w.Code)
	}

	if w.Body.String() != "pong" {
		t.Errorf("Expected body to be 'pong' got %s", w.Body.String())
	}
}

func Test_UsersPostRoute(t *testing.T) {
	db := setupDB()
	defer db.Close()

	router := setupRouter()
	w := httptest.NewRecorder()
	var newUserStr = []byte(`{"name":"Darth Vader", "email": "darth@vader.com"}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(newUserStr))
	router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Expected code 201 got %d", w.Code)
	}

	user := models.User{}
	json.Unmarshal([]byte(w.Body.String()), &user)

	var userFromDb models.User
	db.First(&userFromDb)

	if userFromDb.ID != user.ID {
		t.Error("Expected created and unmarshaled users to be the same")
	}

	if !strings.Contains(w.Body.String(), "\"name\":\"Darth Vader\",\"email\":\"darth@vader.com\"") {
		t.Errorf("Expected body to be json got %s", w.Body.String())
	}

	db.Delete(&user)
}
