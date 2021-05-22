package models

import "time"

// Recovery содержит данные, необходимые для восстановления пароля.
type Restore struct {
	Expired       time.Time `bson:"expired" json:"expired"`           // время убиения токена
	NextAttemptAt time.Time `bson:"next_attempt" json:"next_attempt"` // время в которое разрешается следующая попытка запроса
	Token         string    `bson:"token" json:"token"`               // секретный токен, содержащийся в ссылке восстановления
	Send          bool      `bson:"send" json:"send"`                 // признак успешности того, что письмо отправлено
	Status        string    `bson:"status" json:"status"`             // статус обработки восстановления
	Method        string    `bson:"method" json:"method"`             // метод восстановления enum ["phone", "email"]
	Phone         string    `bson:"phone" json:"phone"`               // выбранный телефон для восстановления
	Email         string    `bson:"email" json:"email"`               // выбранный email для восстановления
	Tries         uint8     `bson:"tries"`                            // Количество попыток для ввода смс
}
