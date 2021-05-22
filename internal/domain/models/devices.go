package models

type Device struct {
	ID string `bson:"id" json:"id" validate:"required"`
	OS string `bson:"signal" json:"signal" validate:"required"`
}
