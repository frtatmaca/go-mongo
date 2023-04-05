package mongoclient

import (
	"context"
	"github.com/firat.atmaca/go-mongo/modules/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var app config.AppTools

type IClient interface {
	Collection(collectionName string) *mongo.Collection
	Disconnect(ctx context.Context) error
	GracefulShutdown(graceTime time.Duration)
}

type Client struct {
	Connection *mongo.Client
}

func NewClient(URI string) (IClient, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		app.ErrorLogger.Panicln(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		app.ErrorLogger.Fatalln(err)
	}

	return &Client{Connection: client}, nil
}

func (c *Client) Collection(collectionName string) *mongo.Collection {
	var collection = c.Connection.Database("go_mongo").Collection(collectionName)

	return collection
}

func (c *Client) Disconnect(ctx context.Context) error {
	err := c.Connection.Disconnect(ctx)
	return err
}

func (c *Client) GracefulShutdown(graceTime time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), graceTime)
	defer func() {
		cancel()
		_ = c.Connection.Disconnect(ctx)
	}()
}
