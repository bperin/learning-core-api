package pubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

type Service struct {
	client *pubsub.Client
	topic  *pubsub.Topic
}

func NewService(ctx context.Context, projectID, topicID string) (*Service, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}

	topic := client.Topic(topicID)

	return &Service{
		client: client,
		topic:  topic,
	}, nil
}

func (s *Service) Close() error {
	return s.client.Close()
}

func (s *Service) Publish(ctx context.Context, msg []byte) error {
	result := s.topic.Publish(ctx, &pubsub.Message{
		Data: msg,
	})
	_, err := result.Get(ctx)
	return err
}
