package myapp

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndexPathHandler(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body)

	assert.Equal(t, "Hello World", string(data))
}

func TestBarPathHandler_WithoutName(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body)

	assert.Equal(t, "Hello World!", string(data))
}

func TestBarPathHandler_WithName(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar?name=woong", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body)

	assert.Equal(t, "Hello woong!", string(data))
}

func TestFooHandler_WithoutJson(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/foo", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestFooHandler_WithJson(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/foo",
		strings.NewReader(`{"first_name":"woong", "last_name":"kim", "email":"woong@gmail.com"}"`))

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)

	user := new(User)
	err := json.NewDecoder(res.Body).Decode(user)
	assert.Nil(t, err)
	assert.Equal(t, "woong", user.FirstName)
	assert.Equal(t, "kim", user.LastName)
}
