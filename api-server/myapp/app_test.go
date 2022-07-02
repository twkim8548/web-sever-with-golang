package myapp

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	assert.New(t)

	ts := httptest.NewServer(MakeNewHandler())
	defer ts.Close()

	response, err := http.Get(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	data, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, string(data), "Hello World")
}

func TestUsers(t *testing.T) {
	assert.New(t)

	ts := httptest.NewServer(MakeNewHandler())
	defer ts.Close()

	response, err := http.Get(ts.URL + "/users")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	data, _ := ioutil.ReadAll(response.Body)
	assert.Contains(t, string(data), "Get UserInfo")
}

func TestGetUserInfo(t *testing.T) {
	assert.New(t)

	ts := httptest.NewServer(MakeNewHandler())
	defer ts.Close()

	response, err := http.Get(ts.URL + "/users/89")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	data, _ := ioutil.ReadAll(response.Body)
	assert.Contains(t, string(data), "No User Id : 89")

	response, err = http.Get(ts.URL + "/users/56")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	data, _ = ioutil.ReadAll(response.Body)
	assert.Contains(t, string(data), "No User Id : 56")
}

func TestCreateUser(t *testing.T) {
	assert.New(t)

	ts := httptest.NewServer(MakeNewHandler())
	defer ts.Close()

	response, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"taewoong", "last_name":"kim", "email" : "rlaxodnd95@naver.com"}`))

	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusCreated)

	user := new(User)
	err = json.NewDecoder(response.Body).Decode(user)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, user.ID)

	id := user.ID
	response, err = http.Get(ts.URL + "/users/" + strconv.Itoa(id))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	user2 := new(User)
	err = json.NewDecoder(response.Body).Decode(user2)
	assert.NoError(t, err)

	assert.Equal(t, user.ID, user2.ID)
	assert.Equal(t, user.FirstName, user2.FirstName)
}
