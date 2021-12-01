package mongo

import (
	"calendar.com/pkg/domain/entity"
	"calendar.com/pkg/storage"
	"context"
	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
	"time"
)

type eventPut struct {
	ID          uuid.UUID `bson:"_id,omitempty"`
	Title       string
	Description string
	Timezone    string
	Time        *time.Time
	Duration    time.Duration
	User        uuid.UUID
	Notes       []entity.Note
}

type Client struct {
	*storage.MongoClient
}

type EventRepository interface {
	Create(put *entity.Event) error
}

func (c *Client) Create(e *entity.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := c.Client.Database(viper.GetString("db_client")).Collection("event").InsertOne(ctx, e)
	return err
}
