package api

import "github.com/gorilla/mux"

func Route() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/cards/open", OpenByCard).Methods("GET")
	r.HandleFunc("/passwords/open", OpenByPassword).Methods("GET")
	return r
}
