package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/edmilsonmedeiross/go-order-generator/internal/infra/mysql"
)

func main() {

	dsn := "orders_user:orders_local@tcp(127.0.0.1:3306)/orders_db"

	db, err := mysql.NewConnection(context.Background(), dsn)
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}

	defer db.Close()

	log.Println("starting server")

	http.HandleFunc("/health", HealthHandelr)

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Println(err)
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
