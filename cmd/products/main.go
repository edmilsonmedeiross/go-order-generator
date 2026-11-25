package main

import (
	"context"
	"log"
	"net/http"

	"github.com/edmilsonmedeiross/go-order-generator/cmd/products/internal/handler"
	"github.com/edmilsonmedeiross/go-order-generator/cmd/products/internal/infra/mysql"
	"github.com/edmilsonmedeiross/go-order-generator/cmd/products/internal/repository"
	"github.com/edmilsonmedeiross/go-order-generator/cmd/products/internal/usecase"
	"github.com/go-playground/validator/v10"
)

func main() {

	dsn := "orders_user:orders_local@tcp(127.0.0.1:3306)/products_db"

	db, err := mysql.NewConnection(context.Background(), dsn)
	if err != nil {
		log.Println("error connecting to database", err)
		return
	}
	defer db.Close()

	validator := validator.New()

	productRepository := repository.NewProductRepository(db)
	productUseCase := usecase.NewProductUseCase(productRepository)
	productHandler := handler.NewProductHandler(productUseCase, validator)

	http.HandleFunc("POST /v1/products", productHandler.CreateProductHandler)
	http.HandleFunc("GET /v1/products", productHandler.ListProductsHandler)
	http.HandleFunc("GET /v1/products/{id}", productHandler.GetProductHandler)

	log.Println("starting server")
	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Println("error starting server", err)
	}
}
