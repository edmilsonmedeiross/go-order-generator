package usecase

import (
	"context"
	"uuid"

	"github.com/edmilsonmedeiross/go-order-generator/internal/domain"
)

type OrderRepository interface {
	Save(ctx context.Context, order *domain.Order) error
	GetProductInfo(ctx context.Context, productID string) (domain.ProductInfo, error)
}

type OrderUseCase struct {
	repo OrderRepository
}

type CreateOrderInput struct {
	CustomerID string
	Name       string
	Items      []CreateOrderItemInput
}

type CreateOrderItemInput struct {
	ProductID  string
	Quantity   int
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, input CreateOrderInput) error {
	order := &domain.Order{
		ID: 				generateID(),
		Status:     domain.OrderStatusCreated,
		CustomerID: input.CustomerID,
		Name:       input.Name,
	}

	for _, item := range input.Items {
		productInfo, err := uc.repo.GetProductInfo(ctx, item.ProductID)
		if err != nil {
			return err
		}
		order.Items = append(order.Items, domain.OrderItem{
			ID: 			 		generateID(),
			OrderID: 	 		order.ID,
			ProductID:  	item.ProductID,
			Quantity:   	item.Quantity,
			PriceCents: 	productInfo.Price,
			ProductName: 	productInfo.Name,
		})
	}

	if err := order.Validate(); err != nil {
		return err
	}


	return uc.repo.Save(ctx, order)
}

func NewOrderUseCase(repo OrderRepository) *OrderUseCase {
	return &OrderUseCase{
		repo: repo,
	}
}

func generateID() string {
	return uuid.New().String()
}