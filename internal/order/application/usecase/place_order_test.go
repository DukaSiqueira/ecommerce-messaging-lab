package usecase

import (
	"context"
	"testing"

	"github.com/DukaSiqueira/ecommerce-messaging-lab/internal/order/domain"
)

type fakeEventPublisher struct {
	receivedEvent domain.OrderPlaced
}

func (fake *fakeEventPublisher) Publish(
	ctx context.Context,
	event domain.OrderPlaced,
) error {
	fake.receivedEvent = event
	return nil
}

func TestPlaceOrderExecuteWithFakePublisher(t *testing.T) {
	fake := &fakeEventPublisher{}

	useCase := NewPlaceOrder(fake)

	input := PlaceOrderInput{
		OrderID:    "order-01",
		CustomerID: "customer-01",
		Items: []PlaceOrderItemInput{
			{
				ProductID: "product-01",
				Quantity:  1,
			},
		},
	}

	err := useCase.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("execute retornou um erro inesperado: %v", err)
	}

	if fake.receivedEvent.OrderID != input.OrderID {
		t.Errorf(
			"OrderID recebido: %q; esperado: %q",
			fake.receivedEvent.OrderID,
			input.OrderID,
		)
	}

	expectedEventID := "order-placed:order-01"
	if fake.receivedEvent.EventID != expectedEventID {
		t.Errorf(
			"EventID recebido: %q; esperado: %q",
			fake.receivedEvent.EventID,
			expectedEventID,
		)
	}

	if fake.receivedEvent.CustomerID != input.CustomerID {
		t.Errorf(
			"CustomerID recebido: %q; esperado: %q",
			fake.receivedEvent.CustomerID,
			input.CustomerID,
		)
	}

	if len(fake.receivedEvent.Items) != len(input.Items) {
		t.Fatalf(
			"quantidade de itens recebida: %d; esperada: %d",
			len(fake.receivedEvent.Items),
			len(input.Items),
		)
	}

	receivedItem := fake.receivedEvent.Items[0]
	expectedItem := input.Items[0]

	if receivedItem.ProductID != expectedItem.ProductID {
		t.Errorf(
			"ProductID recebido: %q; esperado: %q",
			receivedItem.ProductID,
			expectedItem.ProductID,
		)
	}

	if receivedItem.Quantity != expectedItem.Quantity {
		t.Errorf(
			"Quantity recebida: %d; esperada: %d",
			receivedItem.Quantity,
			expectedItem.Quantity,
		)
	}

	if fake.receivedEvent.OccurredAt.IsZero() {
		t.Errorf("OccurredAt não foi preenchido")
	}
}
