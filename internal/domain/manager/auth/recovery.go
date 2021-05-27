package auth

import (
	"context"
	"strings"
	"time"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"github.com/JetBrainer/sso/internal/domain"
	"github.com/JetBrainer/sso/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	restoreTTL               = time.Second * 120 // время жизни токена на восстановление
	emailTokenLen            = 32                // длина токена восстановления для email
	defaultVerifySpamPenalty = time.Minute * 60
	TriesLimit               = 3 // Кол-во попыток при вводе СМС
)

// RecoveryManager управляет процедурой восстановления доступа к аккаунту.
type RecoveryManager struct {
	db          drivers.DataStore
	users       *Users
	spamPenalty time.Duration
	isTesting   bool
}

func (m *Manager) RecoveryManager() *RecoveryManager {
	return &RecoveryManager{
		db:          m.db,
		users:       m.Users(),
		isTesting:   m.isTesting,
		spamPenalty: defaultVerifySpamPenalty,
	}
}

func (rm *RecoveryManager) Testing() {
	rm.isTesting = true
}

func (rm *RecoveryManager) WithSpamPenalty(t time.Duration) *RecoveryManager {
	rm.spamPenalty = t
	return rm
}

// RecoveryByEmail по указанным данным находит пользователя и начинает для него процедуру
// восстановления доступа по адресу почты.
func (rm *RecoveryManager) RecoveryByEmail(ctx context.Context, email string) error {
	if email == "" {
		return ErrEmailNotSpecified
	}

	user, err := rm.users.ByEmail(email)
	if err != nil {
		return ErrEmailNotLinkedToAccount
	}

	token, err := rm.NewTokenForPhone()
	if err != nil {
		return err
	}

	return rm.db.RestoreByEmailNew(ctx, user.ID, email, token, time.Now().In(time.UTC).Add(restoreTTL))
}

// NewTokenForEmail генерирует новый токен.
func (rm *RecoveryManager) NewTokenForEmail() (string, error) {
	return domain.GenerateRandomString(emailTokenLen)
}

// NewTokenForPhone генерирует новый токен.
func (rm *RecoveryManager) NewTokenForPhone() (string, error) {
	if rm.isTesting {
		return strings.Repeat("1", domain.PhoneTokenLen), nil
	}

	return domain.GenerateRandomNumbers(domain.PhoneTokenLen)
}


func (rm *RecoveryManager) FindNewRecoveryRequests(ctx context.Context, c chan<- models.User) {
	rm.db.RestoreFindNew(ctx, c)
}

func (rm *RecoveryManager) RestoreSendNotificationSuccessfully(ctx context.Context, tdid primitive.ObjectID) error {
	return rm.db.RestoreSendNotificationSuccessfully(ctx, tdid)
}