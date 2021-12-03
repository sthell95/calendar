package mongodb

import (
	"context"

	"calendar.com/pkg/logger"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"calendar.com/pkg/domain/entity"
)

const userTable = "users"

type user struct {
	ID       string `bson:"_id,omitempty"`
	Login    string
	Password string
	Timezone string
}

type UserClient struct {
	*mongo.Client
}

func (c *UserClient) FindById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var model *user
	filer := bson.M{"_id": id.String()}

	res := c.Client.Database(viper.GetString("mongo_db")).Collection(userTable).FindOne(ctx, filer)
	err := res.Decode(model)
	if err != nil {
		return nil, err
	}

	u := modelToUser(model)
	return u, nil
}

func (c *UserClient) FindOneByLogin(ctx context.Context, login string) (*entity.User, error) {
	var model *user
	filter := bson.M{"login": login}

	res := c.Client.Database(viper.GetString("mongo_db")).Collection(userTable).FindOne(ctx, filter)
	err := res.Decode(model)
	if err != nil {
		return nil, err
	}

	u := modelToUser(model)
	return u, nil
}

func modelToUser(e *user) *entity.User {
	id, err := uuid.Parse(e.ID)
	if err != nil {
		logger.NewLogger().Write("warning", "Could not parse user id", "mongo-user")
	}

	return &entity.User{
		ID:       id,
		Login:    e.Login,
		Password: e.Password,
		Timezone: e.Timezone,
	}
}

func NewUserRepository(c *mongo.Client) *UserClient {
	return &UserClient{c}
}
