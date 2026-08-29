package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	projectID := os.Getenv("PUBSUB_PROJECT_ID")
	emulatorHost := os.Getenv("PUBSUB_EMULATOR_HOST")

	if projectID == "" {
		log.Fatal("PUBSUB_PROJECT_ID não está configurado")
	}

	if emulatorHost == "" {
		log.Fatal("PUBSUB_EMULATOR_HOST não está configurado")
	}

	fmt.Println("projectID configurado:", projectID)
	fmt.Println("emulatorHost configurado:", emulatorHost)

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("não foi possível criar o cliente Pub/Sub: %v", err)
	}
	defer client.Close()

	topicID := "order-events"
	topicName := fmt.Sprintf(
		"projects/%s/topics/%s",
		projectID,
		topicID,
	)

	topicRequest := &pubsubpb.Topic{
		Name: topicName,
	}
	fmt.Println("tópico a ser criado:", topicRequest.GetName())

	topicCtx, topicCancel := context.WithTimeout(ctx, 5*time.Second)
	defer topicCancel()

	createdTopic, err := client.TopicAdminClient.CreateTopic(
		topicCtx,
		topicRequest,
	)

	if err == nil {
		fmt.Println("tópico criado com sucesso:", createdTopic.GetName())
	} else if status.Code(err) == codes.AlreadyExists {
		fmt.Println("tópico já configurado:", topicRequest.GetName())
	} else {
		log.Printf("não foi possível criar o tópico: %v", err)
		return
	}

	// --------------------------------------------------------------------

	subscriptionsIDs := []string{
		"inventory-order-events-sub",
		"notification-order-events-sub",
	}

	for _, subscriptionID := range subscriptionsIDs {
		err = ensureSubscription(
			ctx,
			client,
			projectID,
			topicName,
			subscriptionID,
		)

		if err != nil {
			log.Printf("erro ao criar subscription: %v", err)
			return
		}
	}

	fmt.Println("configuração do Google Pub/Sub concluída")
}

func ensureSubscription(
	ctx context.Context,
	client *pubsub.Client,
	projectID string,
	topicName string,
	subscriptionID string,
) error {
	subscriptionName := fmt.Sprintf(
		"projects/%s/subscriptions/%s",
		projectID,
		subscriptionID,
	)

	subscriptionRequest := &pubsubpb.Subscription{
		Name:  subscriptionName,
		Topic: topicName,
	}
	fmt.Println("subscription a ser criada:", subscriptionRequest.GetName())
	fmt.Println("tópico da subscription:", subscriptionRequest.GetTopic())

	subscriptionCtx, subscriptionCancel := context.WithTimeout(ctx, 5*time.Second)
	defer subscriptionCancel()

	subscriptionCreated, err := client.SubscriptionAdminClient.CreateSubscription(
		subscriptionCtx,
		subscriptionRequest,
	)

	if err == nil {
		fmt.Println("subscription criada com sucesso:", subscriptionCreated.GetName())
		fmt.Println("subscription do tópico:", subscriptionCreated.GetTopic())
	} else if status.Code(err) == codes.AlreadyExists {
		fmt.Println("subscription já configurada:", subscriptionName)
		fmt.Println("subscription do tópico:", topicName)
	} else {
		return err
	}

	return nil
}
