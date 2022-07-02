package main

import (
	"net/http"
	"web-sever-with-golang/api-server/myapp"
)

func main() {
	http.ListenAndServe(":3000", myapp.MakeNewHandler())
}
