package auth

import (
	"github.com/go-chi/jwtauth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserClaims модель данных пользователя для хранения в JWT токене.
type UserClaims struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	FirstName    string             `bson:"firstname,omitempty" json:"firstname,omitempty"`
	LastName     string             `bson:"lastname,omitempty" json:"lastname,omitempty"`
	Email        string             `bson:"email" json:"email"`
	PrimaryPhone string             `bson:"phone" json:"phone"`
	Phones       []string           `bson:"phones" json:"phones"`
}

type Authenticator struct {
	jwtKey []byte
}

func NewAuthenticator(jwtKey []byte) *Authenticator {
	return &Authenticator{
		jwtKey: jwtKey,
	}
}

// JWT key
func (a *Authenticator) JWTKey() []byte {
	return a.jwtKey
}

func (a *Authenticator) TokenAuth() *jwtauth.JWTAuth {
	return jwtauth.New("HS256", a.jwtKey, nil)
}
