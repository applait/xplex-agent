package rest

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Start bootstraps the REST API
func Start() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", homeHandler).Methods("GET")
	callbackHandler(r.PathPrefix("/callback").Subrouter())

	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
