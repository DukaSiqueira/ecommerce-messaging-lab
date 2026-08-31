package googlepubsub

import (
	"context"

	"cloud.google.com/go/pubsub/v2"
	"github.com/DukaSiqueira/ecommerce-messaging-lab/internal/order/domain"
)

type GooglePubSubPublisher struct {
	publisher *pubsub.Publisher
}

func NewGooglePubSubPublisher(publisher *pubsub.Publisher) *GooglePubSubPublisher {
	return &GooglePubSubPublisher{
		publisher: publisher,
	}
}

func (p *GooglePubSubPublisher) Publish(ctx context.Context, event domain.OrderPlaced) error {
	return nil
}
