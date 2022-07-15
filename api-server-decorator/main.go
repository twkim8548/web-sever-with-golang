package main

import (
	"log"
	"net/http"
	"time"
	"web-sever-with-golang/api-server-decorator/decohandler"
	"web-sever-with-golang/api-server-decorator/myapp"
)

func logger(w http.ResponseWriter, r *http.Request, h http.Handler) {
	start := time.Now()
	log.Println("[LOGGER1] Started")
	h.ServeHTTP(w, r)
	log.Println("[LOGGER1] Completed time:", time.Since(start).Milliseconds())
}

func logger2(w http.ResponseWriter, r *http.Request, h http.Handler) {
	start := time.Now()
	log.Println("[LOGGER2] Started")
	h.ServeHTTP(w, r)
	log.Println("[LOGGER2] Completed time:", time.Since(start).Milliseconds())
}

func NewHandler() http.Handler {
	h := myapp.NewHandler()
	h = decohandler.NewDecoHandler(h, logger)
	h = decohandler.NewDecoHandler(h, logger2)
	return h
}

func main() {

	mux := NewHandler()

	http.ListenAndServe(":3000", mux)
}
