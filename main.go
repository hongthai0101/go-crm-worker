package main

import (
	"context"
	"crm-worker-go/config"
	"crm-worker-go/datasources"
	"crm-worker-go/subscriptions"
	"crm-worker-go/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ctx         context.Context
	mongoClient *mongo.Client
	cancel      context.CancelFunc
)

func init() {
	config.LoadENV()
	utils.InitializeLogger()

	mongoClient, ctx, cancel, _ = datasources.ConnectDB()
}

func main() {
	defer datasources.Close(mongoClient, ctx, cancel)

	subscriptions.Boot()
}
