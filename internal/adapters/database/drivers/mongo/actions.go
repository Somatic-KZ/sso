package mongo

import (
	"context"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	errors2 "github.com/JetBrainer/sso/internal/domain/errors"
	"github.com/JetBrainer/sso/internal/domain/models"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ActionsRepository struct {
	collection *mongo.Collection
}

func (a ActionsRepository) Create(ctx context.Context, action *models.Action) error {
	if action == nil {
		return errors.WithMessage(drivers.ErrEmptyStruct, "cannot create new action")
	}

	actionDocument := bson.D{
		{Key: "title", Value: action.Title},
		{Key: "type", Value: action.Type},
	}
	if _, err := a.collection.InsertOne(ctx, actionDocument); err != nil {
		return errors.Wrap(err, "attempted to create action, got")
	}

	return nil
}

func (a ActionsRepository) ByID(ctx context.Context, id models.PolymorphicID) (*models.Action, error) {
	objID, err := id.ToObjectID()
	if err != nil {
		return nil, errors.WithMessage(errors2.ErrInvalidID, "cannot continue action search")
	}

	action := new(models.Action)
	filter := bson.D{{"_id", objID}}
	err = a.collection.FindOne(ctx, filter).Decode(action)

	switch err {
	case nil:
		return action, nil
	case mongo.ErrNoDocuments:
		return nil, errors.WithMessage(errors2.ErrDoesNotExist, "cannot find action")
	default:
		return nil, errors.Wrap(err, "attempted to find action, got")
	}
}

func (a ActionsRepository) ByType(ctx context.Context, actionType string) (*models.Action, error) {
	action := new(models.Action)
	filter := bson.D{{"type", actionType}}
	err := a.collection.FindOne(ctx, filter).Decode(action)

	switch err {
	case nil:
		return action, nil
	case mongo.ErrNoDocuments:
		return nil, errors.WithMessage(errors2.ErrDoesNotExist, "cannot find action")
	default:
		return nil, errors.Wrap(err, "attempted to find action, got")
	}
}

func (a ActionsRepository) All(ctx context.Context) ([]models.Action, error) {
	actions := make([]models.Action, 0)
	filter := bson.D{}

	cur, err := a.collection.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return actions, nil
		}

		return nil, errors.Wrap(err, "attempted to find actions, got")
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &actions); err != nil {
		return nil, errors.Wrap(err, "could not map actions from datastore")
	}

	return actions, nil
}

func (a ActionsRepository) Update(ctx context.Context, action *models.Action) error {
	if action == nil {
		return errors.WithMessage(drivers.ErrEmptyStruct, "cannot update action")
	}

	objID, err := action.ID.ToObjectID()
	if err != nil {
		return errors.WithMessage(errors2.ErrInvalidID, "cannot continue action search")
	}

	filter := bson.D{{"_id", objID}}
	update := bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "title", Value: action.Title},
			{Key: "type", Value: action.Type},
		},
	}}

	result, err := a.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return errors.WithMessage(errors2.ErrDoesNotExist, "cannot find action")
		default:
			return errors.Wrap(err, "attempted to update action, got")
		}
	}

	if result.ModifiedCount == 0 {
		return errors.WithMessage(errors2.ErrDoesNotExist, "cannot find action")
	}

	return nil
}

func (a ActionsRepository) DeleteByID(ctx context.Context, id models.PolymorphicID) error {
	objID, err := id.ToObjectID()
	if err != nil {
		return errors.WithMessage(errors2.ErrInvalidID, "cannot continue action search")
	}

	filter := bson.D{{"_id", objID}}

	result, err := a.collection.DeleteOne(ctx, filter)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return errors.WithMessage(errors2.ErrDoesNotExist, "cannot find action")
		default:
			return errors.Wrap(err, "attempted to delete action, got")
		}
	}

	if result.DeletedCount == 0 {
		return errors.WithMessage(errors2.ErrDoesNotExist, "cannot find action")
	}

	return nil
}

func (a ActionsRepository) DeleteByType(ctx context.Context, actionType string) error {
	filter := bson.D{{"type", actionType}}

	result, err := a.collection.DeleteOne(ctx, filter)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return errors.WithMessage(errors2.ErrDoesNotExist, "cannot find action")
		default:
			return errors.Wrap(err, "attempted to delete action, got")
		}
	}

	if result.DeletedCount == 0 {
		return errors.WithMessage(errors2.ErrDoesNotExist, "cannot find action")
	}

	return nil
}
