package gcloud

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Constants"
	"time"
)

func CreateTopicOnCloudPubSub(projectID, topicID string) (*pubsub.Topic, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)

	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}

	topic, err := client.CreateTopic(ctx, topicID)

	if err != nil {
		return nil, fmt.Errorf("CreateTopic: %v", err)
	}

	return topic, nil
}

func CreatePullSubscriberForTopicOnCloudPubSub(projectID, subID string, topicID string) (*pubsub.Subscription, error) {

	// topic of type https://godoc.org/cloud.google.com/go/pubsub#Topic
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)

	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}

	topic := client.Topic(topicID)

	sub, err := client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: Constants.MESSAGE_ACK_DEALINE_IN_SECS * time.Second,
		EnableMessageOrdering: Constants.ENABLE_MESSAGE_ORDERING,
	})

	if err != nil {
		return nil, fmt.Errorf("CreateSubscription: %v", err)
	}

	return sub, nil
}

func CreatePushSubscriberForTopicOnCloudPubSub(projectID, subID string, topicID string, endpoint string) (*pubsub.Subscription, error) {
	// projectID := "my-project-id"
	// subID := "my-sub"
	// topic of type https://godoc.org/cloud.google.com/go/pubsub#Topic
	// endpoint := "https://my-test-project.appspot.com/push"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}

	topic := client.Topic(topicID)

	sub, err := client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: Constants.MESSAGE_ACK_DEALINE_IN_SECS * time.Second,
		PushConfig:  pubsub.PushConfig{Endpoint: endpoint},
		EnableMessageOrdering: Constants.ENABLE_MESSAGE_ORDERING,
	})

	if err != nil {
		return nil, fmt.Errorf("CreateSubscription: %v", err)
	}

	return sub, nil
}

func PublishMessageToCloudRunTopic(data []byte, projectID, topicID string) (string, error) {

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return "", fmt.Errorf("pubsub.NewClient: %v", err)
	}

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: data,
		Attributes: map[string]string{
			"origin":   "golang",
			"username": "gcp",
		},
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("Get: %v", err)
	}

	return id, nil
}