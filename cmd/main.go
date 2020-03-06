package main

import (
	"fmt"
	"log"
	"net/http"

	"crud/pkg/config"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/maps", config.MapsPost).Methods("POST")
	router.HandleFunc("/maps", config.MapsGet).Methods("GET")
	router.HandleFunc("/maps", config.MapsPut).Methods("PUT")
	router.HandleFunc("/maps", config.MapsDelete).Methods("DELETE")
	router.HandleFunc("/maps", config.Test).Methods("PATCH")
	http.Handle("/", router)
	fmt.Println("Terhubung dengan port 1612")
	log.Fatal(http.ListenAndServe(":1612", router))
}
