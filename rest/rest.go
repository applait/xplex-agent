package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Start bootstraps the REST API with root namespace handlers
func Start() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", homeHandler).Methods("GET")
	callbackHandler(r.PathPrefix("/callback").Subrouter())
	statsHandler(r.PathPrefix("/stats").Subrouter())

	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
