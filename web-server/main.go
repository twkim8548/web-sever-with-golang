package main

import (
	"net/http"
	"web-sever-with-golang/web-server/myapp"
)

func main() {
	http.ListenAndServe(":3000", myapp.NewHttpHandler())
}
