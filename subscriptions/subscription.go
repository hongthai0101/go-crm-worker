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
	"sync/atomic"
	"time"
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

func (s *Subscription) pullMessages(subscription config.SubscriptionConfigItem) error {
	if client == nil {
		projectID := config.GCSConfig["projectId"]
		ctx = context.Background()
		client, err = pubsub.NewClient(ctx, projectID)
		if err != nil {
			fmt.Errorf("pubsub.NewClient: %v", err)
		}
	}

	defer client.Close()

	sub := client.Subscription(subscription.Key)
	// Must set ReceiveSettings.Synchronous to false (or leave as default) to enable
	// concurrency pulling of messages. Otherwise, NumGoroutines will be set to 1.
	sub.ReceiveSettings.Synchronous = false
	// NumGoroutines determines the number of goroutines sub.Receive will spawn to pull
	// messages.
	sub.ReceiveSettings.NumGoroutines = 16
	// MaxOutstandingMessages limits the number of concurrent handlers of messages.
	// In this case, up to 8 unacked messages can be handled concurrently.
	// Note, even in synchronous mode, messages pulled in a batch can still be handled
	// concurrently.
	sub.ReceiveSettings.MaxOutstandingMessages = 8

	// Receive messages for 10 seconds, which simplifies testing. Comment this out in
	// production, since `Receive` should be used as a long running operation.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var received uint32
	// Receive blocks until the context is cancelled or an error occurs.
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		atomic.AddUint32(&received, 1)

		var isDone bool
		switch subscription.Action {
		case "ExportCrm":
			var payload types.PayloadMessageExport
			if err := json.Unmarshal(msg.Data, &payload); err != nil {
				utils.Logger.Error(err)
			}
			isDone = s.exportService.ExportSaleOpp(payload)
			break
		case "OrderCreated":
			var payload types.RequestMessageOrder
			if err := json.Unmarshal(msg.Data, &payload); err != nil {
				utils.Logger.Error(err)
			}
			isDone = s.saleService.ExecuteMessage(payload, msg.Attributes["source"])
			break
		case "OrderDisbursed":
			var payload types.RequestMessageOrder
			if err := json.Unmarshal(msg.Data, &payload); err != nil {
				utils.Logger.Error(err)
			}
			isDone = s.saleService.ExecuteMessage(payload, msg.Attributes["source"])
			break
		}
		if isDone {
			msg.Ack()
		} else {
			msg.Nack()
		}

	})
	fmt.Printf("Received %d messages\n", received)

	return nil
}

func (s *Subscription) Boot() {
	//for _, item := range config.SubscriptionConfig {
	//	if err = pullMessages(item); err != nil {
	//		utils.Logger.Error(err)
	//	}
	//	println(item.Key)
	//}

	payload := types.RequestMessageOrder{
		Order: types.RequestOrder{
			CustomerName: "Long Vu Dai",
			Email:        "0984536485@gmail.com",
			Phone:        "0984536485",
			AssetType:    "KHC",
			Detail:       "khong co chi",
			Days:         "30",
			Bill:         0,
			Id:           "123",
			CreatedBy:    "0c33cbf8-6212-4454-8e6d-99807b9c3f1d",
			CustomerId:   "0c33cbf8-6212-4454-8e6d-99807b9c3f1d",
		},
		Metadata: nil,
		Images:   []interface{}{"a"},
	}
	s.saleService.ExecuteMessage(payload, "MOBILE")
}
