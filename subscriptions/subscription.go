package subscriptions

import (
	"cloud.google.com/go/pubsub"
	"context"
	"crm-worker-go/config"
	"crm-worker-go/services"
	"crm-worker-go/types"
	"crm-worker-go/utils"
	"encoding/json"
	"fmt"
)

type Subscription struct {
	saleService   *services.SaleService
	exportService *services.ExportService
}

func NewSubscription(service *services.Service) *Subscription {
	return &Subscription{saleService: service.SaleService, exportService: service.ExportService}
}

var (
	client *pubsub.Client
	err    error
	ctx    context.Context
)

func (s *Subscription) pullMessages(subscription *config.SubscriptionConfigItem) error {
	if client == nil {
		projectID := config.GetConfig().GCSConfig.ProjectId
		ctx = context.Background()
		client, err = pubsub.NewClient(ctx, projectID)
		if err != nil {
			fmt.Errorf("pubsub.NewClient: %v", err)
		}
	}

	defer func(client *pubsub.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	sub := client.Subscription(subscription.Key)
	// MaxOutstandingMessages is the maximum number of unprocessed messages the
	// subscriber client will pull from the server before pausing. This also configures
	// the maximum number of concurrent handlers for received messages.
	//
	// For more information, see https://cloud.google.com/pubsub/docs/pull#streamingpull_dealing_with_large_backlogs_of_small_messages.
	sub.ReceiveSettings.MaxOutstandingMessages = 100
	// MaxOutstandingBytes is the maximum size of unprocessed messages,
	// that the subscriber client will pull from the server before pausing.
	sub.ReceiveSettings.MaxOutstandingBytes = 1e8
	// Receive blocks until the context is cancelled or an error occurs.

	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		utils.Logger.Info("Data: ", string(msg.Data), "Attributes: ", msg.Attributes)

		var isDone bool
		switch subscription.Action {
		case "ExportCrm":
			var payload = unmarshalMessage[types.PayloadMessageExport](msg)
			isDone = s.exportService.ExportSaleOpp(payload)
			break
		case "OrderCreated":
			var payload = unmarshalMessage[types.RequestMessageOrder](msg)
			isDone = s.saleService.ExecuteMessage(payload, msg.Attributes["source"])
			break
		case "OrderDisbursed":
			var payload = unmarshalMessage[types.RequestMessageOrder](msg)
			isDone = s.saleService.ExecuteMessage(payload, msg.Attributes["source"])
			break
		}
		println(isDone)
		msg.Ack()
		//if isDone {
		//	msg.Ack()
		//} else {
		//	msg.Nack()
		//}
	})

	if err != nil {
		return err
	}

	fmt.Printf("Received 1 messages of subscription %v \n", subscription.Key)

	return nil
}

func (s *Subscription) Boot() {
	//for _, item := range config.GetConfig().SubscriptionConfig {
	//	if err = s.pullMessages(item); err != nil {
	//		utils.Logger.Error(err)
	//	}
	//}
}

func unmarshalMessage[T interface{}](msg *pubsub.Message) T {
	var payload T
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		utils.Logger.Error(err)
	}
	return payload
}
