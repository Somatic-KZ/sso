package mongo

import (
	"context"
	"log"
	"time"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"github.com/JetBrainer/sso/internal/domain/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RestoreUserByToken возвращает пользователя по токену восстановления.
func (m *Mongo) RestoreUserByToken(ctx context.Context, token string) (*models.User, error) {
	if token == "" {
		return nil, drivers.ErrTokenNotSpec
	}

	var u *models.User

	filter := bson.D{
		{Key: "restore.token", Value: token},
	}

	err := m.DB.Collection(CollectionUsers).FindOne(ctx, filter).Decode(&u)
	switch err {
	case nil:
		return u, nil
	case mongo.ErrNoDocuments:
		return nil, drivers.ErrTokenNotFound
	default:
		return nil, err
	}
}

// RestoreFindExpiredAndUpdate находит просроченные восстановления
// и убирает из них токены.
func (m *Mongo) RestoreFindExpiredAndUpdate(ctx context.Context, c chan<- models.User) {
	collection := m.DB.Collection(CollectionUsers)

	filter := bson.D{
		{Key: "restore.expired",
			Value: bson.D{{Key: "$lte", Value: time.Now().In(time.UTC)}},
		},
		{Key: "restore.token",
			Value: bson.D{{Key: "$ne", Value: ""}},
		},
	}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "restore.token", Value: ""},
			},
		},
	}

	go func() {
		for {
			var u models.User

			singleRes := collection.FindOneAndUpdate(ctx, filter, update)
			err := singleRes.Decode(&u)

			switch err {
			case nil:
				c <- u
			case mongo.ErrNoDocuments:
				return
			default:
				log.Printf("[ERROR] mongo unhandeled error in RestoreFindExpiredAndUpdate(): %v\n", err)
				return
			}
		}
	}()
}

// RestoreFindNew находит все новые восстановления паролей для пользователя.
// Используется для того, чтобы послать им уведомления.
func (m *Mongo) RestoreFindNew(ctx context.Context, c chan<- models.User) {
	collection := m.DB.Collection(CollectionUsers)

	filter := bson.D{
		{Key: "restore.token",
			Value: bson.D{{Key: "$ne", Value: ""}},
		},
		{Key: "restore.send", Value: false},
		{Key: "restore.status", Value: "new"},
	}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "restore.status", Value: "on_check"},
			},
		},
	}

	go func() {
		for {
			var u models.User

			singleRes := collection.FindOneAndUpdate(ctx, filter, update)
			err := singleRes.Decode(&u)

			switch err {
			case nil:
				c <- u
			case mongo.ErrNoDocuments:
				return
			default:
				log.Printf("[ERROR] mongo unhandeled error in RestoreFindNew(): %v\n", err)
				return
			}
		}
	}()
}

// RestoreSendNotificationSuccessfully проставляет признак успешности отправки пользователю
// уведомления о восстановлении пароля.
func (m *Mongo) RestoreSendNotificationSuccessfully(ctx context.Context, tdid primitive.ObjectID) error {
	if tdid.IsZero() {
		return drivers.ErrUserIDNotSpec
	}

	filter := bson.D{
		{Key: "_id", Value: tdid},
	}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "restore.send", Value: true},
				{Key: "restore.status", Value: "finish"},
			},
		},
	}

	_, err := m.DB.Collection(CollectionUsers).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) VerifySendNotificationSuccessfully(ctx context.Context, tdid primitive.ObjectID) error {
	if tdid.IsZero() {
		return drivers.ErrUserIDNotSpec
	}

	filter := bson.D{
		{Key: "_id", Value: tdid},
	}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "verify.send", Value: true},
				{Key: "verify.status", Value: "finish"},
			},
		},
	}

	_, err := m.DB.Collection(CollectionUsers).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) RestoreByEmailNew(ctx context.Context, tdid primitive.ObjectID, email, token string, expiredAt time.Time) error {
	if email == "" {
		return drivers.ErrUserEmailNotSpec
	}

	filter := bson.D{
		{Key: "_id", Value: tdid},
	}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "restore.expired", Value: expiredAt},
				{Key: "restore.token", Value: token},
				{Key: "restore.send", Value: false},
				{Key: "restore.status", Value: "new"},
				{Key: "restore.method", Value: "email"},
				{Key: "restore.email", Value: email},
			},
		},
	}

	result, err := m.DB.Collection(CollectionUsers).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return drivers.ErrUserDoesNotExist
	}

	return nil
}

// RestoreByPhoneNew сохраняет данные для начала процедуры восстановления пароля
// через отсылку sms.
func (m *Mongo) RestoreByPhoneNew(ctx context.Context, tdid primitive.ObjectID, phone, token string, expiredAt, nextAttemptAt time.Time) error {
	if phone == "" {
		return drivers.ErrUserPhoneNotSpec
	}

	filter := bson.D{
		{Key: "_id", Value: tdid},
	}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "restore.expired", Value: expiredAt},
				{Key: "restore.token", Value: token},
				{Key: "restore.send", Value: false},
				{Key: "restore.status", Value: "new"},
				{Key: "restore.method", Value: "phone"},
				{Key: "restore.phone", Value: phone},
				{Key: "restore.next_attempt", Value: nextAttemptAt},
				{Key: "restore.tries", Value: 0},
			},
		},
	}

	result, err := m.DB.Collection(CollectionUsers).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return drivers.ErrUserDoesNotExist
	}

	return nil
}

// RestoreClean завершает процедуру восстановления пароля.
// Используется для удачных и истехших по времени случаев.
func (m *Mongo) RestoreClean(ctx context.Context, tdid primitive.ObjectID) error {
	if tdid.IsZero() {
		return drivers.ErrUserIDNotSpec
	}

	filter := bson.D{
		{Key: "_id", Value: tdid},
	}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "restore.token", Value: ""},
				{Key: "restore.send", Value: false},
				{Key: "restore.tries", Value: 0},
			},
		},
	}

	result, err := m.DB.Collection(CollectionUsers).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return drivers.ErrUserDoesNotExist
	}

	return nil
}

// RestoreIncrementTries Добавляет одну неудачную попытку ввода СМС
func (m *Mongo) RestoreIncrementTries(ctx context.Context, tdid primitive.ObjectID) error {
	if tdid.IsZero() {
		return drivers.ErrUserIDNotSpec
	}

	filter := bson.D{
		{Key: "_id", Value: tdid},
	}
	update := bson.D{
		{Key: "$inc",
			Value: bson.D{
				{Key: "restore.tries", Value: 1},
			},
		},
	}

	result, err := m.DB.Collection(CollectionUsers).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return drivers.ErrUserDoesNotExist
	}

	return nil
}

func (m *Mongo) RestoreUpdate(ctx context.Context, user *models.User) error {
	if user == nil {
		return drivers.ErrEmptyUserStruct
	}

	if user.Restore == nil {
		return errors.Wrap(drivers.ErrEmptyStruct, "restore")
	}

	filter := bson.D{
		{Key: "_id", Value: user.ID},
	}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "restore", Value: user.Restore},
				{Key: "updated", Value: time.Now().In(time.UTC)},
			},
		},
	}

	result, err := m.DB.Collection(CollectionUsers).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return drivers.ErrUserDoesNotExist
	}

	return nil
}
