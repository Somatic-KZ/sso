package auth

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"github.com/JetBrainer/sso/internal/domain"
	"github.com/JetBrainer/sso/internal/domain/errors"
	"github.com/JetBrainer/sso/internal/domain/models"
	"github.com/JetBrainer/sso/internal/ports/http/resources"
	"github.com/JetBrainer/sso/utils"
	"github.com/go-chi/render"
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

// RecoveryByPhone по указанным данным находит пользователя и начинает для него процедуру
// восстановления доступа по номеру телефона.
func (rm *RecoveryManager) RecoveryByPhone(ctx context.Context, phone string) render.Renderer {
	if phone == "" {
		return resources.BadRequest(ErrPhoneNotSpecified)
	}

	normPhone := utils.NormPhoneNum(phone)

	user, err := rm.users.ByAnyPhone(normPhone)
	if err != nil {
		if err == drivers.ErrUserDoesNotExist {
			return resources.ResourceNotFound(ErrUserDoesNotExist)
		}

		log.Printf("[ERROR] Ошибка получения пользователя: %v", err)
		return resources.Internal(ErrInternalError)
	}

	// может быть пользователь многократно вызывает сервис, нельзя допустить
	// многократную отсылку SMS
	if user.Restore != nil && user.Restore.NextAttemptAt.In(time.UTC).After(time.Now().In(time.UTC)) {
		penalty := int(user.Restore.NextAttemptAt.Sub(time.Now()).Seconds())
		return TooManyRequests(errors.ErrTooManyRequests, penalty)
	}

	token, err := rm.NewTokenForPhone()
	if err != nil {
		log.Printf("[ERROR] Ошибка создания токена %v", err)
		return resources.Internal(ErrInternalError)
	}

	nextAttemptAt := time.Now().In(time.UTC).Add(rm.spamPenalty)

	err = rm.db.RestoreByPhoneNew(ctx,
		user.ID,
		phone,
		token,
		time.Now().In(time.UTC).Add(restoreTTL),
		nextAttemptAt,
	)
	if err != nil {
		log.Printf("[ERROR] Ошибка записи токена %v", err)
		return resources.Internal(ErrInternalError)
	}

	return nil
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

// PhoneTokenValidator сравнивает токен восстановления.
func (rm *RecoveryManager) PhoneTokenValidator(phone, token string) error {
	if phone == "" {
		return ErrPhoneNotSpecified
	}

	if token == "" {
		return drivers.ErrTokenNotSpec
	}

	user, err := rm.users.ByAnyPhone(utils.NormPhoneNum(phone))
	if err != nil {
		return errors.ErrTokenDoesNotExist
	}

	if token != user.Restore.Token {
		if user.Restore.Tries == TriesLimit {
			// удаляем токен на восстановления
			if err := rm.db.RestoreClean(context.Background(), user.ID); err != nil {
				log.Printf("[ERROR] RestoreClean() failed: %v", err)
			}

			return errors.ErrTokenTriesExpired
		}

		if user.Restore.Token != "" {
			if err := rm.db.RestoreIncrementTries(context.Background(), user.ID); err != nil {
				log.Printf("[ERROR] IncrementRestoreTries() failed: %v", err)
			}
		}

		return errors.ErrTokenDoesNotExist
	}

	return nil
}

// NewPasswordByPhoneToken меняет пароль по телефонному токену восстановления.
func (rm *RecoveryManager) NewPasswordByPhoneToken(phone, token, newPassword string) error {
	if phone == "" {
		return ErrPhoneNotSpecified
	}

	if token == "" {
		return errors.ErrTokenNotSpec
	}

	phone = utils.NormPhoneNum(phone)

	user, err := rm.users.ByAnyPhone(phone)
	if err != nil {
		return ErrUserDoesNotExist
	}

	if token != user.Restore.Token {
		if user.Restore.Tries == TriesLimit {
			// удаляем токен на восстановления
			if err := rm.db.RestoreClean(context.Background(), user.ID); err != nil {
				log.Printf("[ERROR] RestoreClean() failed: %v", err)
			}

			return errors.ErrTokenTriesExpired
		}

		if user.Restore.Token != "" {
			if err := rm.db.RestoreIncrementTries(context.Background(), user.ID); err != nil {
				log.Printf("[ERROR] IncrementRestoreTries() failed: %v", err)
			}
		}

		return errors.ErrTokenDoesNotExist
	}

	if !user.VerifiedContain(phone) {
		user.Phones = append(user.Phones, phone)
		if err = rm.users.Update(user); err != nil {
			log.Printf("[ERROR] Users.Update() failed: %v", err)
			return ErrInternalError
		}
	}

	if err := rm.users.UpdatePassword(user.ID, newPassword); err != nil {
		return err
	}

	// удаляем токен на восстановления
	if err := rm.db.RestoreClean(context.Background(), user.ID); err != nil {
		log.Printf("[ERROR] RestoreClean() failed: %v", err)
	}

	return nil
}

func (rm *RecoveryManager) FindNewRecoveryRequests(ctx context.Context, c chan<- models.User) {
	rm.db.RestoreFindNew(ctx, c)
}

func (rm *RecoveryManager) RestoreSendNotificationSuccessfully(ctx context.Context, tdid primitive.ObjectID) error {
	return rm.db.RestoreSendNotificationSuccessfully(ctx, tdid)
}

func (rm RecoveryManager) RestoreReset(ctx context.Context, phone string) error {
	user, err := rm.users.ByAnyPhone(utils.NormPhoneNum(phone))
	if err != nil {
		return err
	}

	user.Restore = &models.Restore{
		Status:        domain.TokenStatusFinish,
		Expired:       time.Now().Add(MinusDay),
		NextAttemptAt: time.Now().Add(MinusDay),
	}

	if err := rm.db.RestoreUpdate(ctx, user); err != nil {
		return err
	}

	return nil
}
