package usecase

import "github.com/DukaSiqueira/ecommerce-messaging-lab/internal/order/application/port"

type PlaceOrder struct {
	eventPublisher port.EventPublisher
}

type PlaceOrderInput struct {
	OrderID    string
	CustomerID string
	Items      []PlaceOrderItemInput
}

type PlaceOrderItemInput struct {
	ProductID string
	Quantity  int
}

func NewPlaceOrder(eventPublisher port.EventPublisher) *PlaceOrder {
	return &PlaceOrder{
		eventPublisher: eventPublisher,
	}
}
