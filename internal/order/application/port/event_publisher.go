package port

import (
	"context"

	"github.com/DukaSiqueira/ecommerce-messaging-lab/internal/order/domain"
)

type EventPublisher interface {
	Publish(ctx context.Context, event domain.OrderPlaced) error
}
