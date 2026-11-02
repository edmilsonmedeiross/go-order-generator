package repository

import (
	"context"
	"database/sql"

	"github.com/edmilsonmedeiross/go-order-generator/internal/domain"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Save(ctx context.Context, product *domain.Product) error {
	err := product.Validate()
	if err != nil {
		return err
	}

	prd := &domain.Product{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	}

	_, err = r.db.ExecContext(ctx, saveProductQuery, prd.ID, prd.Name, prd.Price)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) GetProductInfo(ctx context.Context, productID string) (domain.ProductInfo, error) {
	var productInfo domain.ProductInfo
	err := r.db.QueryRowContext(ctx, getProductInfoQuery, productID).Scan(&productInfo.ID, &productInfo.Name, &productInfo.Price)
	if err != nil {
		return domain.ProductInfo{}, err
	}

	return productInfo, nil
}

func (r *ProductRepository) ListProducts(ctx context.Context) ([]domain.ProductInfo, error) {
	rows, err := r.db.QueryContext(ctx, listProductsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.ProductInfo
	for rows.Next() {
		var product domain.ProductInfo
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return products, nil
}
