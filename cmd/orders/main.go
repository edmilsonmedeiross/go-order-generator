package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/edmilsonmedeiross/go-order-generator/cmd/orders/internal/config"
	"github.com/edmilsonmedeiross/go-order-generator/cmd/orders/internal/handler"
	"github.com/edmilsonmedeiross/go-order-generator/cmd/orders/internal/infra/httpclient"
	"github.com/edmilsonmedeiross/go-order-generator/cmd/orders/internal/infra/mysql"
	"github.com/edmilsonmedeiross/go-order-generator/cmd/orders/internal/infra/productclient"
	"github.com/edmilsonmedeiross/go-order-generator/cmd/orders/internal/repository"
	"github.com/edmilsonmedeiross/go-order-generator/cmd/orders/internal/usecase"
	"github.com/go-playground/validator/v10"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	dsn := "orders_user:orders_local@tcp(127.0.0.1:3306)/orders_db"

	db, err := mysql.NewConnection(context.Background(), dsn)
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}

	defer db.Close()

	validator := validator.New()
	orderRepository := repository.NewOrderRepository(db)
	httpClient := httpclient.New(cfg.HTTP)
	productClient := productclient.New(httpClient, cfg.Products)

	orderUseCase := usecase.NewOrderUseCase(orderRepository, productClient)
	orderHandler := handler.NewOrderRequestHandler(orderUseCase, validator)

	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("POST /v1/orders", orderHandler.CreateOrderHandler)

	log.Println("starting server")
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Println(err)
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"message": "ok",
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("error wrighting response", err)
	}
}
