package infra

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

type PubSubService struct {
	client *pubsub.Client
	topic  *pubsub.Topic
}

func NewPubSubService(ctx context.Context, projectID, topicID string) (*PubSubService, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}

	topic := client.Topic(topicID)

	return &PubSubService{
		client: client,
		topic:  topic,
	}, nil
}

func (s *PubSubService) Close() error {
	return s.client.Close()
}

func (s *PubSubService) Publish(ctx context.Context, msg []byte) error {
	result := s.topic.Publish(ctx, &pubsub.Message{
		Data: msg,
	})
	_, err := result.Get(ctx)
	return err
}
