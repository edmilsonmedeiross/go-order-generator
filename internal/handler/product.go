package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/edmilsonmedeiross/go-order-generator/internal/usecase"
	"github.com/go-playground/validator/v10"
)

type productCreator interface {
	CreateProduct(ctx context.Context, product usecase.ProductInput) error
	ListProducts(ctx context.Context) ([]usecase.ProductInfo, error)
}

type ProductRequest struct {
	Name  string `json:"name" validate:"required,min=1"`
	Price int64    `json:"price_cents" validate:"required,gt=0"`
}

type ProductResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price_cents"`
}

type ProductHandler struct {
	usecase 			productCreator
	validator     *validator.Validate
}

func NewProductHandler(productCreator productCreator, validator *validator.Validate) *ProductHandler {
	return &ProductHandler{
		usecase: 		productCreator,
		validator: 	validator,
	}
}

func (ph *ProductHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var product ProductRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := ph.validator.Struct(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prd := usecase.ProductInput{	
		Name:  product.Name,
		Price: product.Price,
	}

	err = ph.usecase.CreateProduct(ctx, prd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ph *ProductHandler) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	products, err := ph.usecase.ListProducts(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prodInfos := make([]ProductResponse, len(products))
	for i, p := range products {
		prodInfos[i] = ProductResponse{
			ID:    p.ID,
			Name:  p.Name,
			Price: p.Price,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(prodInfos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}