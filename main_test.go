package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func Test_UserGetRoute(t *testing.T) {
	db := setupDB()
	defer db.Close()

	userFromDb, _ := models.NewUser("Luke Skywalker", "luke@skywalker.com")
	userFromDb.Create()

	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%d", userFromDb.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected code 200 got %d", w.Code)
	}

	user := models.User{}
	json.Unmarshal([]byte(w.Body.String()), &user)

	if userFromDb.ID != user.ID {
		t.Error("Expected found and unmarshaled users to be the same")
	}

	if !strings.Contains(w.Body.String(), "\"name\":\"Luke Skywalker\",\"email\":\"luke@skywalker.com\"") {
		t.Errorf("Expected body to be json got %s", w.Body.String())
	}

	db.Delete(&userFromDb)
}

func Test_UserGetRouteWhenIdIncorrect(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/luke", nil)
	router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("Expected code 400 got %d", w.Code)
	}

	if !strings.HasPrefix(w.Body.String(), "{\"error\":") {
		t.Error("Expected body to be have error json")
	}
}
