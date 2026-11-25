package usecase

import (
	"context"
	"uuid"

	"github.com/edmilsonmedeiross/go-order-generator/cmd/products/internal/domain"
)

type productRepository interface {
	Save(ctx context.Context, product *domain.Product) error
	ListProducts(ctx context.Context) ([]domain.ProductInfo, error)
	GetProductInfo(ctx context.Context, productID string) (domain.ProductInfo, error)
}

func (uc *ProductUseCase) GetProduct(ctx context.Context, id string) (ProductInfo, error) {
	product, err := uc.repo.GetProductInfo(ctx, id)
	if err != nil {
		return ProductInfo{}, err
	}
	return ProductInfo{ID: product.ID, Name: product.Name, Price: product.Price}, nil
}

type ProductUseCase struct {
	repo productRepository
}

type ProductInput struct {
	Name  string
	Price int64
}

type ProductInfo struct {
	ID    string
	Name  string
	Price int64
}

func NewProductUseCase(repo productRepository) *ProductUseCase {
	return &ProductUseCase{
		repo: repo,
	}
}

func (uc *ProductUseCase) CreateProduct(ctx context.Context, input ProductInput) error {
	product := &domain.Product{
		ID:    generateID(),
		Name:  input.Name,
		Price: input.Price,
	}

	err := product.Validate()
	if err != nil {
		return err
	}

	return uc.repo.Save(ctx, product)
}

func (uc *ProductUseCase) ListProducts(ctx context.Context) ([]ProductInfo, error) {
	products, err := uc.repo.ListProducts(ctx)
	if err != nil {
		return nil, err
	}

	prdInfos := make([]ProductInfo, len(products))
	for i, p := range products {
		prdInfos[i] = ProductInfo{
			ID:    p.ID,
			Name:  p.Name,
			Price: p.Price,
		}
	}

	return prdInfos, nil
}

func generateID() string {
	return uuid.New().String()
}
