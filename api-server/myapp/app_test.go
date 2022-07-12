package myapp

import (
	"encoding/json"
	"fmt"
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

	// 테스트 데이터 생성
	response, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"taewoong", "last_name":"kim", "email" : "rlaxodnd95@naver.com"}`))

	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusCreated)

	response, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"turker", "last_name":"park", "email" : "turker@naver.com"}`))

	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusCreated)

	resp, err := http.Get(ts.URL + "/users")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	users := []*User{}
	err = json.NewDecoder(resp.Body).Decode(&users)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(users))
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

func TestDeleteUser(t *testing.T) {
	assert.New(t)

	ts := httptest.NewServer(MakeNewHandler())
	defer ts.Close()

	req, _ := http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	resp, err := http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 지우긴 했어도 1번 유저는 아직 생성되지 않았음
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(t, string(data), "No User ID:1")

	// 테스트를 위해 1번 유저 먼저 만들기
	response, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"taewoong", "last_name":"kim", "email" : "rlaxodnd95@naver.com"}`))

	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusCreated)

	user := new(User)
	err = json.NewDecoder(response.Body).Decode(user)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, user.ID)

	// 다시 생성한 지우기
	req, _ = http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	resp, err = http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(t, string(data), "Deleted User ID:1")
}

func TestUpdateUser(t *testing.T) {
	assert.New(t)

	ts := httptest.NewServer(MakeNewHandler())
	defer ts.Close()

	req, _ := http.NewRequest("PUT", ts.URL+"/users",
		strings.NewReader(`{"id" : 1, "first_name":"updated", "last_name" : "updated", "email" : "updated@naver.com"}`))
	resp, err := http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(t, string(data), "No User ID:1")

	// 테스트를 위해 1번 유저 먼저 만들기
	response, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"taewoong", "last_name":"kim", "email" : "rlaxodnd95@naver.com"}`))

	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusCreated)

	user := new(User)
	err = json.NewDecoder(response.Body).Decode(user)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, user.ID)

	// 다시 업데이트 하기
	updateStr := fmt.Sprintf(`{"id" : %d, "first_name":"json"}`, user.ID)

	req, _ = http.NewRequest("PUT", ts.URL+"/users",
		strings.NewReader(updateStr))
	resp, err = http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	updateUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(updateUser)

	assert.NoError(t, err)
	assert.Equal(t, updateUser.ID, user.ID)
	assert.Equal(t, updateUser.FirstName, "json")
	assert.Equal(t, updateUser.LastName, user.LastName)
	assert.Equal(t, updateUser.Email, user.Email)

}
