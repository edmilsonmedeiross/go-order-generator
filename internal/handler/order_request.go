package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/edmilsonmedeiross/go-order-generator/internal/usecase"
	"github.com/go-playground/validator/v10"
)

type OrderCreator interface {
	CreateOrder(ctx context.Context, input usecase.CreateOrderInput) error
}
type OrderRequest struct {
	usecase OrderCreator
}

type CreateOrderRequest struct {
	CustomerID string `json:"customer_id" validate:"required"`
	Name       string `json:"name" validate:"required"`
	Items      []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type CreateOrderItemRequest struct {
	ProductID  	string `json:"product_id" validate:"required"`
	Quantity   	int    `json:"quantity" validate:"required,gt=0"`
}


func (or *OrderRequest) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	order, err := parseCreateOrderRequest(r)
	ctx := r.Context()
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = or.usecase.CreateOrder(ctx, order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func NewOrderRequestHandler(useCase OrderCreator) *OrderRequest {
	return &OrderRequest{
		usecase: useCase,
	}
}

func parseCreateOrderRequest(r *http.Request) (usecase.CreateOrderInput, error) {
	var req CreateOrderRequest
	validate := validator.New()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return usecase.CreateOrderInput{}, err
	}

	if err := validate.Struct(req); err != nil {
		return usecase.CreateOrderInput{}, err
	}

	order := usecase.CreateOrderInput{
		CustomerID: req.CustomerID,
		Name:       req.Name,
	}
	
	for _, item := range req.Items {
		order.Items = append(order.Items, usecase.CreateOrderItemInput{
			ProductID:  	item.ProductID,
			Quantity:   	item.Quantity,
		})
	}

	return order, nil
}