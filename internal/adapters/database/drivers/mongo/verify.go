package mongo

import (
	"context"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *Mongo) VerifyToken(ctx context.Context, token string) error {
	err := m.DB.Collection(CollectionUsers).FindOne(ctx, bson.M{"token": token}).Err()
	switch err {
	case mongo.ErrNoDocuments:
		return drivers.ErrTokenNotFound
	case nil:
		return nil
	default:
		return err
	}
}