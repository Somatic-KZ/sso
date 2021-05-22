package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Receiver - получатель, привязанный к юзеру, до которого будет осуществляться доставка
type Receiver struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	FirstName       string             `bson:"firstname" json:"firstname"`
	LastName        string             `bson:"lastname" json:"lastname"`
	Email           string             `bson:"email" json:"email"`
	PrimaryPhone    string             `bson:"phone" json:"phone"`
	AdditionalPhone string             `bson:"additionalPhone,omitempty" json:"additionalPhone,omitempty"`
	Address         *ReceiverAddress   `bson:"address,omitempty" json:"address,omitempty"`
	IsDefault       bool               `bson:"isDefault" json:"isDefault"`
	IsOrganization  bool               `bson:"isOrganization" json:"isOrganization"`
	Organization    *Organization      `bson:"organization,omitempty" json:"organization,omitempty"`
}

// UpdatedReceiver - обновленный получатель
type UpdatedReceiver struct {
	ID              primitive.ObjectID `json:"id"`
	FirstName       string             `json:"firstname" validate:"required"`
	LastName        string             `json:"lastname" validate:"required"`
	Email           string             `json:"email" validate:"required,email"`
	PrimaryPhone    string             `json:"phone" validate:"required,is_phone"`
	AdditionalPhone *string            `json:"additionalPhone,omitempty"`
	Address         *ReceiverAddress   `json:"address,omitempty"`
	IsDefault       *bool              `json:"isDefault,omitempty"`
	IsOrganization  *bool              `json:"isOrganization,omitempty"`
	Organization    *Organization      `json:"organization,omitempty"`
}

// ReceiverAddress - адрес получателя
type ReceiverAddress struct {
	Region    ReceiverRegion `bson:"region" json:"region,omitempty"`
	City      string         `bson:"city" json:"city,omitempty"`
	Street    string         `bson:"street" json:"street,omitempty"`
	House     string         `bson:"house" json:"house,omitempty"`
	Floor     string         `bson:"floor" json:"floor,omitempty"`
	Apartment string         `bson:"apartment" json:"apartment,omitempty"`
	Zipcode   string         `bson:"zipcode" json:"zipcode,omitempty"`
	Geo       AddressGeo     `bson:"geo" json:"geo,omitempty"`
}

// ReceiverRegion - данные региона получателя
type ReceiverRegion struct {
	Code string `bson:"code,omitempty" json:"code,omitempty"`
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	ID   int    `bson:"id,omitempty" json:"id,omitempty"`
}

type ReceiverResponse struct {
	FirstName       string             `bson:"firstname" json:"firstname"`
	LastName        string             `bson:"lastname" json:"lastname"`
	Email           string             `bson:"email" json:"email"`
	PrimaryPhone    string             `bson:"phone" json:"phone"`
	AdditionalPhone string             `bson:"additionalPhone,omitempty" json:"additionalPhone,omitempty"`
	Organization    *Organization      `bson:"organization,omitempty" json:"organization,omitempty"`
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	IsDefault       bool               `bson:"isDefault" json:"isDefault"`
	IsOrganization  bool               `bson:"isOrganization" json:"isOrganization"`
}

type ReceiverAddressResponse struct {
	Region    ReceiverRegion     `bson:"region" json:"region,omitempty"`
	City      string             `bson:"city" json:"city,omitempty"`
	Street    string             `bson:"street" json:"street,omitempty"`
	House     string             `bson:"house" json:"house,omitempty"`
	Floor     string             `bson:"floor" json:"floor,omitempty"`
	Apartment string             `bson:"apartment" json:"apartment,omitempty"`
	Zipcode   string             `bson:"zipcode" json:"zipcode,omitempty"`
	Geo       AddressGeo         `bson:"geo" json:"geo,omitempty"`
	ID        primitive.ObjectID `bson:"_id" json:"id"`
}
