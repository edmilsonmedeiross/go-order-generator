package usecase

import (
	"context"
	"uuid"

	"github.com/edmilsonmedeiross/go-order-generator/cmd/orders/internal/domain"
)

type OrderRepository interface {
	Save(ctx context.Context, order *domain.Order) error
}

type ProductClient interface {
	GetProduct(ctx context.Context, id string) (ProductInfo, error)
}

type OrderUseCase struct {
	repo          OrderRepository
	productClient ProductClient
}

type CreateOrderInput struct {
	CustomerID string
	Name       string
	Items      []CreateOrderItemInput
}

type CreateOrderItemInput struct {
	ProductID string
	Quantity  int
}

type ProductInfo struct {
	ID    string "json:\"id\""
	Name  string "json:\"name\""
	Price int64  "json:\"price_cents\""
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, input CreateOrderInput) error {
	order := &domain.Order{
		ID:         generateID(),
		Status:     domain.OrderStatusCreated,
		CustomerID: input.CustomerID,
		Name:       input.Name,
	}

	for _, item := range input.Items {

		productInfo, err := uc.productClient.GetProduct(ctx, item.ProductID)
		if err != nil {
			return err
		}

		order.Items = append(order.Items, domain.OrderItem{
			ID:          generateID(),
			OrderID:     order.ID,
			ProductID:   productInfo.ID,
			Quantity:    item.Quantity,
			PriceCents:  productInfo.Price,
			ProductName: productInfo.Name,
		})
	}

	if err := order.Validate(); err != nil {
		return err
	}

	return uc.repo.Save(ctx, order)
}

func NewOrderUseCase(repo OrderRepository, productClient ProductClient) *OrderUseCase {
	return &OrderUseCase{
		repo:          repo,
		productClient: productClient,
	}
}

func generateID() string {
	return uuid.New().String()
}
