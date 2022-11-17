package datasources

import (
	"context"
	"crm-worker-go/config"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var (
	MongoClient *mongo.Client
)

func ConnectDB() (*mongo.Client, context.Context, context.CancelFunc, error) {
	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			log.Print(evt.Command)
		},
	}

	dbHost := config.GetConfig().DB.Host
	dbPort := config.GetConfig().DB.Port
	dbUser := config.GetConfig().DB.User
	dbPass := config.GetConfig().DB.Pass
	dbName := config.GetConfig().DB.Name
	dbDebug := config.GetConfig().DB.Debug

	mongoURI := fmt.Sprintf("mongodb://%v:%v@%v:%v/?authSource=%v", dbUser, dbPass, dbHost, dbPort, dbName)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	option := options.Client().ApplyURI(mongoURI)
	if dbDebug {
		option.SetMonitor(cmdMonitor)
	}
	client, err := mongo.Connect(ctx, option)
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	MongoClient = client
	return client, ctx, cancel, err
}

func Close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Disconnect To Database Errors %v", err)
		}
	}()
}
