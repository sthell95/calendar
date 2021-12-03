package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"

	"calendar.com/pkg/domain/entity"
)

const eventTable = "events"

type event struct {
	ID          string `bson:"_id,omitempty"`
	Title       string
	Description string
	Timezone    string
	Time        *time.Time
	Duration    time.Duration
	User        string
	Notes       []string
}

type EventClient struct {
	*mongo.Client
}

func (c *EventClient) Create(ctx context.Context, e *entity.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	eventModel := eventToModel(e)
	id, err := c.Client.Database(viper.GetString("mongo_db")).Collection(eventTable).InsertOne(ctx, eventModel)
	e.ID, err = uuid.Parse(id.InsertedID.(string))

	return err
}

func (c *EventClient) Update(ctx context.Context, e *entity.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	eventModel := bson.D{{"$set", eventToModel(e)}}
	_, err := c.Client.Database(viper.GetString("mongo_db")).Collection(eventTable).UpdateByID(ctx, e.ID.String(), eventModel)

	return err
}

func (c *EventClient) Delete(ctx context.Context, e *entity.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": e.ID.String(), "user": e.User.ID.String()}
	_, err := c.Client.Database(viper.GetString("mongo_db")).Collection(eventTable).DeleteOne(ctx, filter)

	return err
}

func (c *EventClient) FindOneById(id *uuid.UUID) (*entity.Event, error) {
	return nil, errors.New("not implemented")
}

func eventToModel(e *entity.Event) *event {
	notes := []string{}
	for _, note := range e.Notes {
		notes = append(notes, note.Note)
	}

	return &event{
		ID:          e.ID.String(),
		Title:       e.Title,
		Description: e.Description,
		Timezone:    e.Timezone,
		Time:        e.Time,
		Duration:    e.Duration,
		User:        e.User.ID.String(),
		Notes:       notes,
	}
}

func modelToEvent(e *event) *entity.Event {
	eventID, _ := uuid.Parse(e.ID)
	userId, _ := uuid.Parse(e.User)

	var notes []entity.Note
	for i, note := range e.Notes {
		notes[i].Note = note
	}

	return &entity.Event{
		ID:          eventID,
		Title:       e.Title,
		Description: e.Description,
		Timezone:    e.Timezone,
		Time:        e.Time,
		Duration:    e.Duration,
		User: entity.User{
			ID: userId,
		},
		Notes: notes,
	}
}

func NewEventRepository(c *mongo.Client) *EventClient {
	return &EventClient{c}
}
