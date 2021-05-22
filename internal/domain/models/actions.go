package models


type Action struct {
	ID    PolymorphicID `bson:"_id" json:"id"`
	Title string               `bson:"title" json:"title" validate:"required"`
	Type  string               `bson:"type" json:"type" validate:"required"`
}
