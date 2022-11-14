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

var (
	client *pubsub.Client
	err    error
	ctx    context.Context
)

func pullMessages(subscription config.SubscriptionConfigItem) error {
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

	// Receive messages for 10 seconds, which simplifies testing.
	// Comment this out in production, since `Receive` should
	// be used as a long running operation.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var received uint32
	// Receive blocks until the context is cancelled or an error occurs.
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		atomic.AddUint32(&received, 1)

		var isDone bool
		switch subscription.Action {
		case "ExportCrm":
			exportService := services.NewExportService()
			var payload types.PayloadMessageExport
			if err := json.Unmarshal(msg.Data, &payload); err != nil {
				utils.Logger.Error(err)
			}
			isDone = exportService.ExportSaleOpp(payload)
			break
		case "OrderCreated":
			saleService := services.NewSaleService()
			var payload types.RequestMessageOrder
			if err := json.Unmarshal(msg.Data, &payload); err != nil {
				utils.Logger.Error(err)
			}
			isDone = saleService.ExecuteMessage(payload, msg.Attributes["source"])
			break
		case "OrderDisbursed":
			saleService := services.NewSaleService()
			var payload types.RequestMessageOrder
			if err := json.Unmarshal(msg.Data, &payload); err != nil {
				utils.Logger.Error(err)
			}
			isDone = saleService.ExecuteMessage(payload, msg.Attributes["source"])
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

func Boot() {
	for _, item := range config.SubscriptionConfig {
		if err = pullMessages(item); err != nil {
			utils.Logger.Error(err)
		}
		println(item.Key)
	}
}
