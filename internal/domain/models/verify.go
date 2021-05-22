package models

import "time"

// Verify содержит данные, необходимые для проведения верификации.
type Verify struct {
	Email         string    `bson:"email" json:"email"`
	Phone         string    `bson:"phone" json:"phone"`
	Token         string    `bson:"token" json:"token"`
	Status        string    `bson:"status" json:"status"` // статус обработки верификации
	Send          bool      `bson:"send" json:"send"`
	Expired       time.Time `bson:"expired" json:"expired"`
	NextAttemptAt time.Time `bson:"next_attempt" json:"next_attempt"`
	Tries         uint8     `bson:"tries"`      // Количество попыток для ввода смс
	Generation    uint8     `bson:"generation"` // Количество сгенерированных смс
}
