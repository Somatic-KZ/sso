package models

import (
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PolymorphicID string

func PolymorphicIDFromString(id string) PolymorphicID {
	return PolymorphicID(id)
}

func (id PolymorphicID) ToObjectID() (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(string(id))
}

func (id PolymorphicID) ToInt64() (int64, error) {
	return strconv.ParseInt(string(id), 10, 64)
}

func (id PolymorphicID) ToInt() (int, error) {
	return strconv.Atoi(string(id))
}