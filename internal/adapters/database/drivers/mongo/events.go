package mongo

import (
	"context"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"github.com/JetBrainer/sso/internal/domain"
	errors2 "github.com/JetBrainer/sso/internal/domain/errors"

	"github.com/JetBrainer/sso/internal/domain/models"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventsRepository struct {
	collection *mongo.Collection
}

func (e EventsRepository) Create(ctx context.Context, event *models.Event) error {
	if event == nil {
		return errors.WithMessage(drivers.ErrEmptyStruct, "cannot create new event")
	}

	objTDID, err := event.TDID.ToObjectID()
	if err != nil {
		return errors.WithMessage(errors2.ErrInvalidID, "cannot create new event with provided tdid")
	}

	eventDocument := bson.D{
		{Key: "tdid", Value: objTDID},
		{Key: "actionType", Value: event.ActionType},
	}
	if event.Verify != nil {
		eventDocument = append(eventDocument, bson.E{Key: "verify", Value: event.Verify})
	}
	if event.ExpiresAt != nil {
		eventDocument = append(eventDocument, bson.E{Key: "expiresAt", Value: event.ExpiresAt})
	}

	if _, err := e.collection.InsertOne(ctx, eventDocument); err != nil {
		return errors.Wrap(err, "attempted to create event, got")
	}

	return nil
}

func (e EventsRepository) ByID(ctx context.Context, id models.PolymorphicID) (*models.Event, error) {
	objID, err := id.ToObjectID()
	if err != nil {
		return nil, errors.WithMessage(errors2.ErrInvalidID, "cannot continue event search")
	}

	event := new(models.Event)
	filter := bson.D{{"_id", objID}}
	err = e.collection.FindOne(ctx, filter).Decode(event)

	switch err {
	case nil:
		return event, nil
	case mongo.ErrNoDocuments:
		return nil, errors.WithMessage(errors2.ErrDoesNotExist, "cannot find event")
	default:
		return nil, errors.Wrap(err, "attempted to find event, got")
	}
}

func (e EventsRepository) ByUserAction(ctx context.Context, tdid models.PolymorphicID, actionType string) (*models.Event, error) {
	objID, err := tdid.ToObjectID()
	if err != nil {
		return nil, errors.WithMessage(errors2.ErrInvalidID, "cannot continue event search for provided tdid")
	}

	event := new(models.Event)
	filter := bson.D{{"tdid", objID}, {"actionType", actionType}}
	err = e.collection.FindOne(ctx, filter).Decode(event)

	switch err {
	case nil:
		return event, nil
	case mongo.ErrNoDocuments:
		return nil, errors.WithMessage(errors2.ErrDoesNotExist, "cannot find event")
	default:
		return nil, errors.Wrap(err, "attempted to find event, got")
	}
}

func (e EventsRepository) All(ctx context.Context) ([]models.Event, error) {
	filter := bson.D{}

	events := make([]models.Event, 0)
	cur, err := e.collection.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return events, nil
		}

		return nil, errors.Wrap(err, "attempted to find events, got")
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &events); err != nil {
		return nil, errors.Wrap(err, "could not map events from datastore")
	}

	return events, nil
}

func (e EventsRepository) VerifyFindNew(ctx context.Context) (*models.Event, error) {
	filter := bson.D{
		{Key: "verify.token",
			Value: bson.D{{Key: "$ne", Value: ""}},
		},
		{Key: "verify.send", Value: false},
		{Key: "verify.status", Value: domain.TokenStatusNew},
	}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "verify.status", Value: domain.TokenStatusOnCheck},
			},
		},
	}

	event := new(models.Event)
	updatedDocStatus := options.After
	opts := &options.FindOneAndUpdateOptions{ReturnDocument: &updatedDocStatus}
	if err := e.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(event); err != nil {
		return nil, err
	}

	return event, nil
}

func (e EventsRepository) Update(ctx context.Context, event *models.Event) error {
	if event == nil {
		return errors.WithMessage(errors2.ErrEmptyStruct, "cannot update event")
	}

	objID, err := event.ID.ToObjectID()
	if err != nil {
		return errors.WithMessage(errors2.ErrInvalidID, "cannot continue event search")
	}

	filter := bson.D{{"_id", objID}}
	eventDocument := bson.D{
		{Key: "actionType", Value: event.ActionType},
	}
	if event.Verify != nil {
		eventDocument = append(eventDocument, bson.E{Key: "verify", Value: event.Verify})
	}
	if event.ExpiresAt != nil {
		eventDocument = append(eventDocument, bson.E{Key: "expiresAt", Value: event.ExpiresAt})
	}
	update := bson.D{{
		Key:   "$set",
		Value: eventDocument,
	}}

	result, err := e.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return errors.WithMessage(errors2.ErrDoesNotExist, "cannot find event")
		default:
			return errors.Wrap(err, "attempted to update event, got")
		}
	}

	if result.ModifiedCount == 0 {
		return errors.WithMessage(errors2.ErrDoesNotExist, "cannot find event")
	}

	return nil
}

func (e EventsRepository) ApproveVerificationSendStatus(ctx context.Context, tdid models.PolymorphicID, actionType string) error {
	objID, err := tdid.ToObjectID()
	if err != nil {
		return errors.WithMessage(errors2.ErrInvalidID, "cannot continue event search for provided tdid")
	}

	filter := bson.D{{"tdid", objID}, {"actionType", actionType}}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "verify.send", Value: true},
			},
		},
	}

	result, err := e.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return errors.WithMessage(errors2.ErrDoesNotExist, "cannot find event")
		default:
			return errors.Wrap(err, "attempted to update event, got")
		}
	}

	if result.ModifiedCount == 0 {
		return errors.WithMessage(errors2.ErrDoesNotExist, "cannot find event")
	}

	return nil
}

func (e EventsRepository) DeleteByID(ctx context.Context, id models.PolymorphicID) error {
	objID, err := id.ToObjectID()
	if err != nil {
		return errors.WithMessage(errors2.ErrInvalidID, "cannot continue event search")
	}

	filter := bson.D{{"_id", objID}}
	result, err := e.collection.DeleteOne(ctx, filter)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return errors.WithMessage(errors2.ErrDoesNotExist, "cannot find event")
		default:
			return errors.Wrap(err, "attempted to delete event, got")
		}
	}

	if result.DeletedCount == 0 {
		return errors.WithMessage(errors2.ErrDoesNotExist, "cannot find event")
	}

	return nil
}
