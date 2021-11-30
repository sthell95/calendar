package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"calendar.com/pkg/logger"
)

type MongoClient struct {
	*mongo.Client
}

func (c *MongoClient) Create(entity interface{}, model Model) error {
	fmt.Sprintln("хуй")
	return nil
}

func (c *MongoClient) Update(entity interface{}, model Model, condition string) error {
	return nil
}

func (c *MongoClient) FindById(entity interface{}, id uuid.UUID) error {
	return nil
}

func (c *MongoClient) FindOneBy(entity interface{}, conditions map[string]interface{}) error {
	return nil
}

func (c *MongoClient) Delete(entity interface{}, model Model, condition string) error {
	return nil
}

func NewDb(ctx context.Context) *MongoClient {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("mongo_url")))
	if err != nil {
		logger.NewLogger().Write(logger.Error, "Mongo url is invalid", "db")
		log.Fatalln(err)
	}

	go func(c *mongo.Client) {
		<-ctx.Done()
		_ = c.Disconnect(ctx)
	}(client)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.NewLogger().Write(logger.Error, "Mongo service has fallen down", "db")
		log.Fatalln(err)
	}

	return &MongoClient{client}
}
