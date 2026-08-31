package usecase

import (
	"context"
	"time"

	"github.com/DukaSiqueira/ecommerce-messaging-lab/internal/order/application/port"
	"github.com/DukaSiqueira/ecommerce-messaging-lab/internal/order/domain"
)

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

func (useCase *PlaceOrder) Execute(
	ctx context.Context,
	input PlaceOrderInput,
) error {
	var eventItems []domain.OrderPlacedItem

	for _, item := range input.Items {
		convertedItem := domain.OrderPlacedItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}

		eventItems = append(eventItems, convertedItem)
	}

	event := domain.OrderPlaced{
		OrderID: input.OrderID,
		EventID: "order-placed:" + input.OrderID,
		Items: eventItems,
		CustomerID: input.CustomerID,
		OccurredAt: time.Now().UTC(),
	}

	return useCase.eventPublisher.Publish(ctx, event)
}
