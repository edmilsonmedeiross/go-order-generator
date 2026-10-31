package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {

	log.Println("starting server")

	http.HandleFunc("/health", HealthHandelr)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func HealthHandelr(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"message": "ok",
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("error wrighting response", err)
	}
}
