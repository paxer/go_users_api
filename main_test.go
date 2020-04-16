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

func Test_UserGetRouteWhenIdNotFound(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/7822345", nil)
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected code 404 got %d", w.Code)
	}

	if w.Body.String() != "{}" {
		t.Errorf("Expected body to be empty JSON, got %s", w.Body.String())
	}
}

func Test_UserUpdateRoute(t *testing.T) {
	db := setupDB()
	defer db.Close()

	user, _ := models.NewUser("Luke Skywalker", "luke@skywalker.com")
	user.Create()

	router := setupRouter()
	w := httptest.NewRecorder()
	var updatedUserJSON = []byte(`{"name":"Darth Vader", "email": "darth@vader.com"}`)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/users/%d", user.ID), bytes.NewBuffer(updatedUserJSON))
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected code 200 got %d", w.Code)
	}

	userResp := models.User{}
	json.Unmarshal([]byte(w.Body.String()), &userResp)

	var userFromDb models.User
	db.First(&userFromDb)

	if userFromDb.Name != "Darth Vader" {
		t.Errorf("Expected Name to be updated but got %s", userFromDb.Name)
	}

	if userFromDb.Email != "darth@vader.com" {
		t.Errorf("Expected Email to be updated but got %s", userFromDb.Email)
	}

	if userResp.ID != userFromDb.ID {
		t.Error("Expected updated and unmarshaled users to be the same")
	}

	if !strings.Contains(w.Body.String(), "\"name\":\"Darth Vader\",\"email\":\"darth@vader.com\"") {
		t.Errorf("Expected body to be json got %s", w.Body.String())
	}

	db.Delete(&user)
}

func Test_UserUpdateRouteWhenIdNotFound(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/users/7822345", nil)
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected code 404 got %d", w.Code)
	}

	if w.Body.String() != "{}" {
		t.Errorf("Expected body to be empty JSON, got %s", w.Body.String())
	}
}

func Test_UserUpdateRouteWhenIdInvalid(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/users/foo", nil)
	router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("Expected code 400 got %d", w.Code)
	}

	if !strings.HasPrefix(w.Body.String(), "{\"error\":") {
		t.Error("Expected body to be have error json")
	}
}

func Test_UserDeleteRouteWhenIdNotFound(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/7822345", nil)
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected code 404 got %d", w.Code)
	}

	if w.Body.String() != "{}" {
		t.Errorf("Expected body to be empty JSON, got %s", w.Body.String())
	}
}

func Test_UserDeleteRouteWhenIdInvalid(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/foo", nil)
	router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("Expected code 400 got %d", w.Code)
	}

	if !strings.HasPrefix(w.Body.String(), "{\"error\":") {
		t.Error("Expected body to be have error json")
	}
}

func Test_UserDeleteRoute(t *testing.T) {
	db := setupDB()
	defer db.Close()

	userFromDb, _ := models.NewUser("Luke Skywalker", "luke@skywalker.com")
	userFromDb.Create()

	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%d", userFromDb.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected code 200 got %d", w.Code)
	}

	var foundUser models.User
	db.First(&foundUser, userFromDb.ID)

	if foundUser.ID != 0 {
		t.Error("Expected user to be deleted from the database")
	}

	if w.Body.String() != "{}" {
		t.Errorf("Expected body to be empty JSON, got %s", w.Body.String())
	}
}

func Test_GetAllUsersRoute(t *testing.T) {
	db := setupDB()
	defer db.Close()

	userFromDb, _ := models.NewUser("Luke Skywalker", "luke@skywalker.com")
	userFromDb.Create()

	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected code 200 got %d", w.Code)
	}

	users := []models.User{}
	json.Unmarshal([]byte(w.Body.String()), &users)

	if len(users) != 1 {
		t.Errorf("Expected one user to be returned got %d", len(users))
	}

	if users[0].ID != userFromDb.ID {
		t.Error("Expected found and unmarshaled users to be the same")
	}

	if !strings.HasPrefix(w.Body.String(), "[{") {
		t.Errorf("Expected body to be json array got %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "\"name\":\"Luke Skywalker\",\"email\":\"luke@skywalker.com\"") {
		t.Errorf("Expected body to be json got %s", w.Body.String())
	}

	db.Delete(&userFromDb)
}
