package main

import (
	"log"
	"net/http"
)

// Root handler
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(string(""))
}
