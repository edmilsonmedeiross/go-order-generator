package repository

import (
	"context"
	"database/sql"

	"github.com/edmilsonmedeiross/go-order-generator/internal/domain"
)

type OrderRepository struct {
	db *sql.DB
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