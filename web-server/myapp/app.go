package myapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type fooHandler struct{}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Hello World")
}

func (f *fooHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	user := new(User)

	err := json.NewDecoder(request.Body).Decode(user)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, "Bad Request: ", err)
		return
	}

	user.CreatedAt = time.Now()

	data, _ := json.Marshal(user)
	writer.Header().Add("content-type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	fmt.Fprint(writer, string(data))
}

func barHandler(writer http.ResponseWriter, request *http.Request) {
	name := request.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(writer, "Hello %s!", name)
}

func NewHttpHandler() http.Handler {
	mux := http.NewServeMux()

	// 라우팅
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/bar", barHandler)
	mux.Handle("/foo", &fooHandler{})

	return mux
}
