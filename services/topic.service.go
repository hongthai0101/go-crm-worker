package services

import (
	"cloud.google.com/go/pubsub"
	"context"
	"crm-worker-go/config"
	"crm-worker-go/utils"
	"encoding/json"
)

type TopicService interface {
	Send(topicName string, body interface{}, attributes map[string]string) bool
}

type topicService struct{}

func NewTopicService() TopicService {
	return &topicService{}
}

func (s *topicService) Send(topicName string, body interface{}, attributes map[string]string) bool {

	utils.Logger.Info(body, attributes)

	ctx := context.Background()
	projectID := config.GCSConfig["projectId"]
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {

	}
	topic := client.Topic(topicName)

	data, _ := json.Marshal(body)
	var msg = &pubsub.Message{
		Data:       data,
		Attributes: attributes,
	}
	if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
		return false
	}
	return true
}
