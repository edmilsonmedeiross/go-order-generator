package domain

import "errors"

type ProductInfo struct {
	ID    string
	Name  string
	Price int64
}

type Product struct {
	ID    string
	Name  string
	Price int64
}

func (p *Product) Validate() error {
	if p.ID == "" {
		return errors.New(ErrInvalidProductID)
	}
	if p.Name == "" {
		return errors.New(ErrInvalidProductName)
	}
	if p.Price <= 0 {
		return errors.New(ErrInvalidProductPrice)
	}

	return nil
}

const (
	ErrInvalidProductID    = "invalid product ID"
	ErrInvalidProductName  = "invalid product name"
	ErrInvalidProductPrice = "invalid product price"
)