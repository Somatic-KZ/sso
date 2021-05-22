package api

import (
	"time"

	"github.com/JetBrainer/sso/internal/domain/models"
)

type ProfileResponse struct {
	Created      time.Time         `bson:"created" json:"created"`
	Updated      time.Time         `bson:"updated" json:"updated"`
	BirthDate    *time.Time        `bson:"birth_date,omitempty" json:"birth_date,omitempty"`
	Roles        []string          `bson:"roles" json:"roles"`
	Permissions  []string          `json:"permissions,omitempty"`
	Receivers    []models.Receiver `bson:"receivers" json:"receivers"`
	FirstName    string            `bson:"firstname,omitempty" json:"firstname"`
	LastName     string            `bson:"lastname,omitempty" json:"lastname"`
	Patronymic   string            `bson:"patronymic,omitempty" json:"patronymic"`
	Email        string            `bson:"email" json:"email"`
	PrimaryPhone string            `bson:"phone" json:"phone"`
	Phones       []string          `bson:"phones" json:"phones"`
	Language     string            `bson:"lang" json:"lang"`
	ID           string            `bson:"_id" json:"tdid"`
	Sex          string            `bson:"sex,omitempty" json:"sex"`
	IIN          int               `bson:"iin,omitempty" json:"iin"`
}

type ProfileUpdateRequest struct {
	FirstName  *string    `bson:"firstname,omitempty" json:"firstname,omitempty"`
	LastName   *string    `bson:"lastname,omitempty" json:"lastname,omitempty"`
	Patronymic *string    `bson:"patronymic,omitempty" json:"patronymic,omitempty"`
	Email      *string    `bson:"email,omitempty" json:"email,omitempty" validate:"omitempty,email,unique_email"`
	Language   *string    `bson:"lang,omitempty" json:"lang,omitempty" validation:"omitempty,oneof='ru en kk'"`
	Sex        *string    `bson:"sex,omitempty" json:"sex,omitempty" validation:"omitempty,oneof='male female unknown'"`
	IIN        *int       `bson:"iin,omitempty" json:"iin,omitempty" validate:"omitempty,iin"`
	BirthDate  *time.Time `bson:"birth_date,omitempty" json:"birth_date,omitempty" validation:"omitempty,datetime=2006-01-02"`
}

type AddPhoneRequest struct {
	Phone string `bson:"phone" json:"phone"`
}

type UpdatePasswordRequest struct {
	Password string `bson:"password" json:"password" validate:"required,gte=6"`
}

type ProfileTypeRequest struct {
	Phone string `json:"phone" validate:"required,is_phone"`
}

type ProfileKindRequest struct {
	Email string `json:"email,omitempty" validate:"omitempty,email"`
	Phone string `json:"phone" validate:"required,is_phone"`
}

type ProfileTypeResponse struct {
	Type string `json:"type"`
}

type AddReceiverRequest struct {
	FirstName       string                  `json:"firstname" validate:"required"`
	LastName        string                  `json:"lastname" validate:"required"`
	Email           string                  `json:"email" validate:"required,email"`
	PrimaryPhone    string                  `json:"phone" validate:"required,is_phone"`
	AdditionalPhone string                  `json:"additionalPhone,omitempty"`
	Address         *models.ReceiverAddress `json:"address,omitempty"`
	IsDefault       bool                    `json:"isDefault,omitempty"`
	IsOrganization  bool                    `json:"isOrganization,omitempty"`
	Organization    *models.Organization    `json:"organization,omitempty"`
}

type UpdateReceiverRequest struct {
	FirstName       string                  `json:"firstname" validate:"required"`
	LastName        string                  `json:"lastname" validate:"required"`
	Email           string                  `json:"email" validate:"required,email"`
	PrimaryPhone    string                  `json:"phone" validate:"required,is_phone"`
	AdditionalPhone *string                 `json:"additionalPhone,omitempty"`
	Address         *models.ReceiverAddress `json:"address,omitempty"`
	IsDefault       *bool                   `json:"isDefault,omitempty"`
	IsOrganization  *bool                   `json:"isOrganization,omitempty"`
	Organization    *models.Organization    `json:"organization,omitempty"`
}

type AddReceiverResponse struct {
	ID string `json:"id"`
}
