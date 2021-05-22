package mongo

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"github.com/JetBrainer/sso/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

)

func (m *Mongo) VerifyPhone(ctx context.Context, tdid primitive.ObjectID, phone, token string, expiredAt, nextAttemptAt time.Time) error {
	if tdid.IsZero() {
		return drivers.ErrUserDoesNotExist
	}

	filter := bson.D{
		{Key: "_id", Value: tdid},
	}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "verify.expired", Value: expiredAt},
				{Key: "verify.token", Value: token},
				{Key: "verify.send", Value: false},
				{Key: "verify.status", Value: "new"},
				{Key: "verify.phone", Value: phone},
				{Key: "verify.next_attempt", Value: nextAttemptAt},
				{Key: "verify.tries", Value: 0},
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

// VerifyClean завершает процедуру верификации телефона.
// Используется для удачных и истехших по времени случаев.
func (m *Mongo) VerifyClean(ctx context.Context, tdid primitive.ObjectID, phone string) error {
	if tdid.IsZero() {
		return drivers.ErrUserIDNotSpec
	}

	filter := bson.D{
		{Key: "_id", Value: tdid},
	}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "verify.token", Value: ""},
				{Key: "verify.send", Value: false},
				{Key: "verify.status", Value: "finish"},
				{Key: "verify.phone", Value: ""},
				{Key: "verify.tries", Value: 0},
				{Key: "phone", Value: phone},
				{Key: "updated", Value: time.Now().In(time.UTC)},
			},
		},
		{Key: "$addToSet", Value: bson.D{{Key: "phones", Value: phone}}},
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

func (m *Mongo) UpdateVerify(ctx context.Context, user *models.User) error {
	if user == nil {
		return drivers.ErrEmptyUserStruct
	}

	filter := bson.D{
		{Key: "_id", Value: user.ID},
	}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "verify", Value: user.Verify},
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

func nextAttempt(generation uint8, spamPenalty time.Duration) time.Time {
	multiplier := time.Duration(math.Pow(float64(generation+1), 2))
	return time.Now().UTC().Add(spamPenalty * multiplier)
}

// VerifyPhoneIncrementGeneration
// Начинает процедуру верификации телефона.
// В случае если пользователь уже запрашивал OTP через СМС,
// то добавляем время до следующей попытки запросить OTP.
func (m *Mongo) VerifyPhoneIncrementGeneration(ctx context.Context, tdid primitive.ObjectID, phone, token string, ttl, spamPenalty time.Duration) error {
	if tdid.IsZero() {
		return drivers.ErrUserIDNotSpec
	}

	var user models.User
	filter := bson.D{{Key: "_id", Value: tdid}}

	err := m.DB.Collection(CollectionUsers).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return err
	}

	verify := user.Verify

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "verify.expired", Value: time.Now().UTC().Add(ttl)},
				{Key: "verify.token", Value: token},
				{Key: "verify.send", Value: false},
				{Key: "verify.status", Value: "new"},
				{Key: "verify.phone", Value: phone},
				{Key: "verify.next_attempt", Value: nextAttempt(verify.Generation, spamPenalty)},
				{Key: "verify.tries", Value: 0},
				{Key: "verify.generation", Value: verify.Generation + 1},
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

// VerifyIncrementTries добавляет неудачную попытку проверки номера через СМС
func (m *Mongo) VerifyIncrementTries(ctx context.Context, tdid primitive.ObjectID) error {
	if tdid.IsZero() {
		return drivers.ErrUserIDNotSpec
	}

	filter := bson.D{
		{Key: "_id", Value: tdid},
	}
	update := bson.D{
		{
			Key: "$inc",
			Value: bson.D{
				{Key: "verify.tries", Value: 1},
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

// VerifyFindNew находит все новые верификации для пользователя.
// Используется для того, чтобы послать им уведомления.
func (m *Mongo) VerifyFindNew(ctx context.Context, c chan<- models.User) {
	collection := m.DB.Collection(CollectionUsers)

	filter := bson.D{
		{Key: "verify.token",
			Value: bson.D{{Key: "$ne", Value: ""}},
		},
		{Key: "verify.send", Value: false},
		{Key: "verify.status", Value: "new"},
	}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "verify.status", Value: "on_check"},
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
				log.Printf("[ERROR] mongo unhandeled error in VerifyFindNew(): %v\n", err)
				return
			}
		}
	}()
}
