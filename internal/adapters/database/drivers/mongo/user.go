package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"github.com/JetBrainer/sso/internal/domain/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CollectionUsers = "users"


// UserCreate создает пользователя и возвращает tdid.
func (m *Mongo) UserCreate(ctx context.Context, user *models.User) (primitive.ObjectID, error) {
	if user == nil {
		return primitive.NilObjectID, drivers.ErrEmptyUserStruct
	}

	if user.Login == "" {
		return primitive.NilObjectID, drivers.ErrUserLoginNotSpec
	}

	result, err := m.DB.Collection(CollectionUsers).InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid, nil
	}

	return primitive.NilObjectID, primitive.ErrInvalidHex
}

// UserByLogin returns the user by the given name or nil.
func (m *Mongo) UserByLogin(ctx context.Context, login string) (*models.User, error) {
	if login == "" {
		return nil, drivers.ErrUserLoginNotSpec
	}

	var user *models.User
	err := m.DB.Collection(CollectionUsers).FindOne(ctx, bson.D{{Key: "login", Value: login}}).Decode(&user)

	switch err {
	case nil:
		return user, nil
	case mongo.ErrNoDocuments:
		return nil, drivers.ErrUserDoesNotExist
	default:
		return nil, err
	}
}

// ByEmail returns the user by the given email or nil.
func (m *Mongo) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user *models.User

	nonCaseSenseMail := fmt.Sprintf("^%s$", email)
	filter := bson.D{
		{Key: "email",
			Value: bson.D{
				{Key: "$regex", Value: nonCaseSenseMail},
				{Key: "$options", Value: "i"},
			},
		},
	}

	err := m.DB.Collection(CollectionUsers).FindOne(ctx, filter).Decode(&user)
	switch err {
	case nil:
		return user, nil
	case mongo.ErrNoDocuments:
		return nil, drivers.ErrUserDoesNotExist
	default:
		return nil, err
	}
}

// ByPhone возвращает пользователя по его номеру телефона.
func (m *Mongo) UserByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user *models.User

	filter := bson.D{
		{Key: "phones",
			Value: bson.D{
				{Key: "$in", Value: []string{phone}},
			},
		},
	}

	err := m.DB.Collection(CollectionUsers).FindOne(ctx, filter).Decode(&user)
	switch err {
	case nil:
		return user, nil
	case mongo.ErrNoDocuments:
		return nil, drivers.ErrUserDoesNotExist
	default:
		return nil, err
	}
}

// UserByPrimaryPhone находит пользователя по первичному
// номеру телефона.
func (m *Mongo) UserByPrimaryPhone(ctx context.Context, phone string) (*models.User, error) {
	var user *models.User

	if phone == "" {
		return nil, drivers.ErrUserPhoneNotSpec
	}

	filter := bson.D{
		{Key: "phone", Value: phone},
	}

	result := m.DB.Collection(CollectionUsers).FindOne(ctx, filter)
	err := result.Decode(&user)
	switch err {
	case nil:
		return user, nil
	case mongo.ErrNoDocuments:
		return nil, drivers.ErrUserDoesNotExist
	default:
		return nil, err
	}
}

func (m *Mongo) UserLastDeviceTokenUpdate(ctx context.Context, tdid string, device *models.Device) error {
	id, err := primitive.ObjectIDFromHex(tdid)
	if err != nil {
		return err
	}
	filter := bson.D{
		{Key: "_id", Value: id},
	}

	update := bson.D{
		{Key: "last_verified.device", Value: device},
	}

	_, err = m.DB.Collection(CollectionUsers).UpdateOne(ctx,filter,update)
	switch err {
	case mongo.ErrNoDocuments:
		return drivers.ErrUserDoesNotExist
	case nil:
		return nil
	default:
		return err
	}
}

func (m *Mongo) LastDeviceToken(ctx context.Context, tdid string) (*models.Device,error) {
	filter := bson.D{{Key: "_id", Value: tdid}}

	user := new(models.User)
	if err := m.DB.Collection(CollectionUsers).FindOne(ctx, filter).Decode(user); err != nil {
		return nil, err
	}

	return user.LastVerified.Device, nil
}

// UsersCount возвращает кол-во пользователей.
func (m *Mongo) UsersCount(ctx context.Context, filters *models.UsersSearchFilters) (int64, error) {
	total, err := m.DB.Collection(CollectionUsers).CountDocuments(ctx, m.searchFilters(filters), &options.CountOptions{})
	if err != nil {
		return 0, err
	}

	return total, nil
}

// UserDelete - удаляет пользователя по его логину.
func (m *Mongo) UserDelete(ctx context.Context, login string) error {
	if login == "" {
		return drivers.ErrUserLoginNotSpec
	}

	result, err := m.DB.Collection(CollectionUsers).DeleteOne(ctx, bson.D{{Key: "login", Value: login}})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return drivers.ErrUserDoesNotExist
	}

	return nil
}

// UserUpdate обновляет структуру models.User по его SignInByLogin.
func (m *Mongo) UserUpdate(ctx context.Context, user *models.User) error {
	if user == nil {
		return drivers.ErrEmptyUserStruct
	}

	if user.ID.IsZero() {
		return drivers.ErrUserIDNotSpec
	}

	filter := bson.D{{Key: "_id", Value: user.ID}}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "firstname", Value: user.FirstName},
				{Key: "lastname", Value: user.LastName},
				{Key: "patronymic", Value: user.Patronymic},
				{Key: "sex", Value: user.Sex},
				{Key: "password", Value: user.Password},
				{Key: "email", Value: user.Email},
				{Key: "enabled", Value: user.Enabled},
				{Key: "phone", Value: user.PrimaryPhone},
				{Key: "phones", Value: user.Phones},
				{Key: "receivers", Value: user.Receivers},
				{Key: "website", Value: user.Website},
				{Key: "updated", Value: time.Now().In(time.UTC)},
				{Key: "roles", Value: user.Roles},
				{Key: "lang", Value: user.Language},
				{Key: "iin", Value: user.IIN},
				{Key: "devices", Value: user.Devices},
				{Key: "bankData", Value: user.BankData},
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

// UserDeleteByTDID удаляет пользователя по его tdid.
func (m *Mongo) UserDeleteByTDID(ctx context.Context, uid primitive.ObjectID) error {
	_, err := m.DB.Collection(CollectionUsers).DeleteOne(ctx, bson.D{{Key: "_id", Value: uid}})
	return err
}

// UserByLogin returns the user by the given tdid or nil.
func (m *Mongo) UserByTDID(ctx context.Context, tdid primitive.ObjectID) (*models.User, error) {
	if tdid.IsZero() {
		return nil, drivers.ErrUserIDNotSpec
	}

	var user *models.User
	err := m.DB.Collection(CollectionUsers).FindOne(ctx, bson.D{{Key: "_id", Value: tdid}}).Decode(&user)

	switch err {
	case nil:
		return user, nil
	case mongo.ErrNoDocuments:
		return nil, drivers.ErrUserDoesNotExist
	default:
		return nil, err
	}
}

func (m *Mongo) UserFullNamesByTDID(ctx context.Context, tdidList []primitive.ObjectID) (map[string]string, error) {
	fullNames := make(map[string]string)

	filter := bson.D{{Key: "_id", Value: bson.D{
		{Key: "$in", Value: tdidList},
	}}}

	cursor, err := m.DB.Collection(CollectionUsers).Find(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		user := new(models.User)
		if err := cursor.Decode(user); err != nil {
			return nil, err
		}

		fullNames[user.ID.Hex()] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}

	return fullNames, nil
}

func (m *Mongo) UserDevicesByTDID(ctx context.Context, tdid primitive.ObjectID) ([]models.Device, error) {
	filter := bson.D{{Key: "_id", Value: tdid}}

	user := new(models.User)
	if err := m.DB.Collection(CollectionUsers).FindOne(ctx, filter).Decode(user); err != nil {
		return nil, err
	}

	return user.Devices, nil
}
func (m *Mongo) searchFilters(filters *models.UsersSearchFilters) bson.D {
	queryFilters := bson.D{}

	if filters != nil {
		if filters.Phone != nil {
			queryFilters = append(queryFilters, bson.E{Key: "phone", Value: primitive.Regex{
				Pattern: *filters.Phone,
				Options: "i",
			}})
		}

		if filters.Email != nil {
			queryFilters = append(queryFilters, bson.E{Key: "email", Value: primitive.Regex{
				Pattern: *filters.Email,
				Options: "i",
			}})
		}

		if filters.Firstname != nil {
			queryFilters = append(queryFilters, bson.E{Key: "firstname", Value: primitive.Regex{
				Pattern: *filters.Firstname,
				Options: "i",
			}})
		}

		if filters.Lastname != nil {
			queryFilters = append(queryFilters, bson.E{Key: "lastname", Value: primitive.Regex{
				Pattern: *filters.Lastname,
				Options: "i",
			}})
		}

		if filters.Patronymic != nil {
			queryFilters = append(queryFilters, bson.E{Key: "patronymic", Value: primitive.Regex{
				Pattern: *filters.Patronymic,
				Options: "i",
			}})
		}
	}

	return queryFilters
}
