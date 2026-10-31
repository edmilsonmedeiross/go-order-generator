package domain

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderTotal(t *testing.T) {
	order := Order{
		ID:         "1",
		CustomerID: "customer1",
		Name:       "Order 1",
		Status:     OrderStatusCreated,
		Items: []OrderItem{
			{
				ProductID:  "product1",
				Quantity:   2,
				PriceCents: 1000,
			},
			{
				ProductID:  "product2",
				Quantity:   1,
				PriceCents: 2000,
			},
		},
	}

	total := order.Total()
	assert.Equal(t, int64(4000), total)
}

func TestOrderValidate(t *testing.T) {
	validOrder := Order{
		ID:         "1",
		CustomerID: "customer1",
		Name:       "Order 1",
		Status:     OrderStatusCreated,
		Items: []OrderItem{
			{
				ProductID:  "product1",
				Quantity:   2,
				PriceCents: 1000,
			},
		},
	}

	err := validOrder.Validate()
	assert.NoError(t, err)
}

func TestOrderValidateItemsInvalid(t *testing.T) {
	invalidOrder := Order{
		ID:         "1",
		CustomerID: "customer1",
		Name:       "Order 1",
		Status:     OrderStatusCreated,
	}

	err := invalidOrder.Validate()
	log.Println("Validation error:", err)
	assert.Error(t, err)
}