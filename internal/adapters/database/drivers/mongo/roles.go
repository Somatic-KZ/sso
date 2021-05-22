package mongo

import (
	"context"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"github.com/JetBrainer/sso/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RolesRepository struct {
	collection *mongo.Collection
}

func (rr *RolesRepository) Roles(ctx context.Context) ([]models.Role, error) {
	cursor, err := rr.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	roles := make([]models.Role, 0)
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}

func (rr *RolesRepository) Create(ctx context.Context, role *models.Role) error {
	if role == nil {
		return drivers.ErrEmptyRoleStruct
	}

	role.ID = primitive.NewObjectID()
	_, err := rr.collection.InsertOne(ctx, role)

	return err
}

func (rr *RolesRepository) RoleByName(ctx context.Context, name string) (*models.Role, error) {
	role := new(models.Role)

	if err := rr.collection.FindOne(ctx, bson.D{{Key: "name", Value: name}}).Decode(role); err != nil {
		return nil, err
	}

	return role, nil
}

func (rr *RolesRepository) Update(ctx context.Context, role *models.Role) error {
	if role == nil {
		return drivers.ErrEmptyRoleStruct
	}

	filter := bson.D{{Key: "name", Value: role.Name}}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "permissions", Value: role.Permissions},
			},
		},
	}

	result, err := rr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return drivers.ErrRoleDoesNotExist
	}

	return nil
}

func (rr *RolesRepository) DeleteByName(ctx context.Context, name string) error {
	_, err := rr.collection.DeleteOne(ctx, bson.D{{Key: "name", Value: name}})
	return err
}
