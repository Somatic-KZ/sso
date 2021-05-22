package models

// Address содержит физические адреса клиентов.
type Address struct {
	Region    string     `bson:"region,omitempty" json:"region,omitempty"`
	City      string     `bson:"city,omitempty" json:"city,omitempty"`
	Street    string     `bson:"street,omitempty" json:"street,omitempty"`
	Corpus    string     `bson:"corpus,omitempty" json:"corpus,omitempty"`
	House     string     `bson:"house,omitempty" json:"house,omitempty"`
	Apartment string     `bson:"apartment,omitempty" json:"apartment,omitempty"`
	Zipcode   int        `bson:"zipcode,omitempty" json:"zipcode,omitempty"`
	Geo       AddressGeo `bson:"geo,omitempty" json:"geo,omitempty"`
}

// AddressGeo содержит данные по гео позиции.
type AddressGeo struct {
	Lat string `bson:"lat" json:"lat"`
	Lng string `bson:"lng" json:"lng"`
}
