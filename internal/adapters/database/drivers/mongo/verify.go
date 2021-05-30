package mongo

import (
	"context"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *Mongo) VerifyToken(ctx context.Context, id string) error {
	tdid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	err = m.DB.Collection(CollectionUsers).FindOne(ctx, bson.M{"_id": tdid}).Err()
	switch err {
	case mongo.ErrNoDocuments:
		return drivers.ErrTokenNotFound
	case nil:
		return nil
	default:
		return err
	}
}