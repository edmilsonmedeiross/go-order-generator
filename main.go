package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/edmilsonmedeiross/go-order-generator/internal/handler"
	"github.com/edmilsonmedeiross/go-order-generator/internal/infra/mysql"
	"github.com/edmilsonmedeiross/go-order-generator/internal/repository"
	"github.com/edmilsonmedeiross/go-order-generator/internal/usecase"
	"github.com/go-playground/validator/v10"
)

func main() {

	dsn := "orders_user:orders_local@tcp(127.0.0.1:3306)/orders_db"

	db, err := mysql.NewConnection(context.Background(), dsn)
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}

	defer db.Close()

	validator := validator.New()

	orderRepository := repository.NewOrderRepository(db)
	orderUseCase := usecase.NewOrderUseCase(orderRepository)
	orderHandler := handler.NewOrderRequestHandler(orderUseCase)
	productRepository := repository.NewProductRepository(db)
	productUseCase := usecase.NewProductUseCase(productRepository)
	productHandler := handler.NewProductHandler(productUseCase, validator)

	http.HandleFunc("/health", HealthHandelr)
	http.HandleFunc("POST /v1/orders", orderHandler.CreateOrderHandler)
	http.HandleFunc("POST /v1/products", productHandler.CreateProductHandler)
	http.HandleFunc("GET /v1/products", productHandler.ListProductsHandler)

	log.Println("starting server")
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
