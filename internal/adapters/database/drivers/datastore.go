package drivers

import (
	"context"
	"time"

	"github.com/JetBrainer/sso/internal/domain/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataStore interface {
	// расширяем функционал datastore здесь
	Name() string
	Close() error
	Connect() error

	UserCreate(ctx context.Context, user *models.User) (primitive.ObjectID, error)
	UserByLogin(ctx context.Context, login string) (*models.User, error)
	UserByTDID(ctx context.Context, tdid primitive.ObjectID) (*models.User, error)
	UserFullNamesByTDID(ctx context.Context, tdidList []primitive.ObjectID) (map[string]string, error)
	UserDevicesByTDID(ctx context.Context, tdid primitive.ObjectID) ([]models.Device, error)
	UserByEmail(ctx context.Context, email string) (*models.User, error)
	UsersCount(ctx context.Context, filters *models.UsersSearchFilters) (int64, error)
	UserDelete(ctx context.Context, login string) error
	UserDeleteByTDID(ctx context.Context, uid primitive.ObjectID) error
	UserUpdate(ctx context.Context, user *models.User) error
	UserByPhone(ctx context.Context, phone string) (*models.User, error)
	UserByPrimaryPhone(ctx context.Context, phone string) (*models.User, error)

	// LastUserDeviceToken обновляет последний использованный девайс токен для отправки уведомлений
	UserLastDeviceTokenUpdate(ctx context.Context, tdid string, device *models.Device) error
	// LastDeviceToken получает последний обновленный токен девайса
	LastDeviceToken(ctx context.Context, tdid string) (*models.Device, error)

	// восстановление доступа
	RestoreByPhoneNew(ctx context.Context, tdid primitive.ObjectID, phone, token string, expiredAt, nextAttemptAt time.Time) error
	RestoreByEmailNew(ctx context.Context, tdid primitive.ObjectID, email, token string, expiredAt time.Time) error
	RestoreClean(ctx context.Context, tdid primitive.ObjectID) error
	RestoreIncrementTries(ctx context.Context, tdid primitive.ObjectID) error
	RestoreSendNotificationSuccessfully(ctx context.Context, tdid primitive.ObjectID) error
	RestoreFindNew(ctx context.Context, c chan<- models.User)
	RestoreFindExpiredAndUpdate(ctx context.Context, c chan<- models.User)
	RestoreUpdate(ctx context.Context, user *models.User) error

	// рутина верификации
	VerifyPhone(ctx context.Context, tdid primitive.ObjectID, phone, token string, expiredAt, nextAttemptAt time.Time) error
	VerifyPhoneIncrementGeneration(ctx context.Context, tdid primitive.ObjectID, phone, token string, expire, nextAttempt time.Duration) error
	VerifyClean(ctx context.Context, tdid primitive.ObjectID, phone string) error
	VerifyIncrementTries(ctx context.Context, tdid primitive.ObjectID) error
	VerifyFindNew(ctx context.Context, c chan<- models.User)
	VerifySendNotificationSuccessfully(ctx context.Context, tdid primitive.ObjectID) error
	UpdateVerify(ctx context.Context, user *models.User) error

	Roles() RolesRepository

	// Проверка на работоспоособность
	Ping() error
}
