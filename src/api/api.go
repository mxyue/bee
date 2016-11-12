package api

import "github.com/gorilla/mux"

func Route() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/cards/open", OpenByCard).Methods("GET")
	return r
}
