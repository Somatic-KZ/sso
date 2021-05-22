package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RolePermissions map[string][]string

type Role struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Permissions RolePermissions    `bson:"permissions" json:"permissions"`
}
