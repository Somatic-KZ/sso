package models

// Organization содержит данные о компании.
type Organization struct {
	Name    string `bson:"name" json:"name"`
	BIN     string `bson:"bin" json:"bin"`
	BIK     string `bson:"bik" json:"bik"`
	IIC     string `bson:"iic" json:"iic"`
	Address string `bson:"address,omitempty" json:"address,omitempty"`
}
