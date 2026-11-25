package domain

import "errors"

type Order struct {
	ID         string
	CustomerID string
	Name       string
	Status     OrderStatus
	Items      []OrderItem
}

type OrderStatus string

const (
	OrderStatusCreated     OrderStatus = "created"
	OrderStatusPaid        OrderStatus = "paid"
	OrderStatusFulfillment OrderStatus = "fulfillment"
	OrderStatusDelivered   OrderStatus = "delivered"
)

type OrderItem struct {
	ID         string
	OrderID    string
	ProductID  string
	Quantity   int
	PriceCents int64
	ProductName string
}

func (o *Order) Total() int64 {
	var total int64
	for _, item := range o.Items {
		total += item.PriceCents * int64(item.Quantity)
	}

	return total
}

func (o *Order) Validate() error {
	if o.ID == "" {
		return errors.New(ErrInvalidOrderID)
	}
	if o.CustomerID == "" {
		return errors.New(ErrInvalidCustomerID)
	}
	if o.Name == "" {
		return errors.New(ErrInvalidOrderName)
	}
	if len(o.Items) == 0 {
		return errors.New(ErrInvalidOrderTotalItems)
	}
	if o.Status != OrderStatusCreated && o.Status != OrderStatusPaid && o.Status != OrderStatusFulfillment && o.Status != OrderStatusDelivered {
		return errors.New(ErrInvalidOrderStatus)
	}

	for _, item := range o.Items {
		if item.PriceCents <= 0 {
			return errors.New(ErrInvalidOrderItemPrice)
		}
		if item.Quantity <= 0 {
			return errors.New(ErrInvalidOrderItemQty)
		}
		if item.ProductID == "" {
			return errors.New(ErrInvalidOrderItemID)
		}
		if item.ProductName == "" {
			return errors.New(ErrInvalidOrderItemName)
		}
	}

	return nil
}

const (
	ErrInvalidOrderID        		= "invalid order ID"
	ErrInvalidCustomerID     		= "invalid customer ID"
	ErrInvalidOrderName      		= "invalid order name"
	ErrInvalidOrderTotal     		= "invalid order total"
	ErrInvalidOrderItemID    		= "invalid order item ID"
	ErrInvalidOrderItemQty   		= "invalid order item quantity"
	ErrInvalidOrderItemPrice 		= "invalid order item price"
	ErrInvalidOrderStatus   		= "invalid order status"
	ErrInvalidOrderTotalItems  	= "invalid order total items"
	ErrInvalidOrderItemName  		= "invalid order item name"
)