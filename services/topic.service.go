package services

import (
	"cloud.google.com/go/pubsub"
	"context"
	"crm-worker-go/config"
	"crm-worker-go/utils"
	"encoding/json"
)

type TopicService struct{}

func NewTopicService() *TopicService {
	return &TopicService{}
}

func (s *TopicService) Send(topicName string, body interface{}, attributes map[string]string) bool {

	utils.Logger.Info(body, attributes)

	ctx := context.Background()
	projectID := config.GetConfig().GCSConfig.ProjectId
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		utils.Logger.Error(err)
		return false
	}
	topic := client.Topic(topicName)

	data, _ := json.Marshal(body)
	var msg = &pubsub.Message{
		Data:       data,
		Attributes: attributes,
	}
	if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
		utils.Logger.Error(err)
		return false
	}
	return true
}
