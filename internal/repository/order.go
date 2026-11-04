package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"uuid"

	"github.com/edmilsonmedeiross/go-order-generator/internal/domain"
)

type OrderRepository struct {
	db *sql.DB
}

type OrderCreatedEvent struct {
	OrderID    string          					`json:"order_id"`
	CustomerID string          					`json:"customer_id"`
	Name       string          					`json:"name"`
	Status     string          					`json:"status"`
	Items      []OrderCreatedItemEvent 	`json:"items"`
	TotalCents int64             					`json:"total_cents"`
}

type OrderCreatedItemEvent struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	PriceCents  int64    `json:"price_cents"`
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Save(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	
	defer tx.Rollback()
	
	_, err = tx.ExecContext(ctx, saveOrderQuery,
		order.ID, order.Status, order.CustomerID, order.Name)
	if err != nil {
		return err
	}

	
	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, saveOrderItemQuery,
			item.ID, item.OrderID, item.ProductID, item.Quantity, item.PriceCents, item.ProductName)
			if err != nil {
				return err
			}
			
		}
		
	orderEvent := &OrderCreatedEvent{
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		Name:       order.Name,
		Status:     string(order.Status),
		Items:      make([]OrderCreatedItemEvent, len(order.Items)),
		TotalCents: order.Total(),
	}

	for i, item := range order.Items {
		orderEvent.Items[i] = OrderCreatedItemEvent{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			PriceCents:  item.PriceCents,
		}
	}

	payload, err := json.Marshal(orderEvent)
	if err != nil {
		return err
	}
	
	_, err = tx.ExecContext(ctx, saveOutboxEventQuery,
		uuid.New().String(), "OrderCreated", order.ID, payload)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *OrderRepository) GetProductInfo(ctx context.Context, productID string) (domain.ProductInfo, error) {
	var productInfo domain.ProductInfo
	err := r.db.QueryRowContext(ctx, getProductInfoQuery, productID).Scan(&productInfo.ID, &productInfo.Name, &productInfo.Price)
	if err != nil {
		return domain.ProductInfo{}, err
	}
	return productInfo, nil
}