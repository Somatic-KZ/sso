package auth

import (
	"strings"
	"time"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"github.com/JetBrainer/sso/internal/domain"
	"github.com/JetBrainer/sso/internal/domain/models"
	"github.com/JetBrainer/sso/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	VerifyTTL              = time.Minute * 5  // время жизни токена на верификацию
	DefaultLongSpamPenalty = time.Minute * 60 // время пенальти при попытках спама восстановления токена без указанного tdid
	DefaultSpamPenalty     = time.Second * 20 // срабатывает при повторных попытках отослать SMS
	MinusDay               = time.Hour * -24
)

// VerificationManager занимается вопросами верификации телефонов и email'ов.
type VerificationManager struct {
	db              drivers.DataStore
	users           *Users
	isTesting       bool
	spamPenalty     time.Duration
	longSpamPenalty time.Duration
}

// UserDataForTokenGeneration содержит данные необходимые для генерации токена
type UserDataForTokenGeneration struct {
	Phone string
	Email string
	TDID  string
}

// VerificationTokenData содержит данные, которые позволяют контролировать процесс верфикации
// смс-токена.
type VerificationTokenData struct {
	User          *models.User
	NormPhone     string
	NextAttemptAt time.Time
}

func (m *Manager) VerificationManager() *VerificationManager {
	return &VerificationManager{
		db:              m.db,
		users:           m.Users(),
		isTesting:       m.isTesting,
		spamPenalty:     DefaultSpamPenalty,
		longSpamPenalty: DefaultLongSpamPenalty,
	}
}

func (vm *VerificationManager) Testing() {
	vm.isTesting = true
}

// WithSpamPenalty устанавливает время пенальти при повторных попытках
// восстановления токена с указанным tdid.
func (vm *VerificationManager) WithSpamPenalty(t time.Duration) *VerificationManager {
	vm.spamPenalty = t
	return vm
}

// NewTokenForPhone генерирует новый токен.
func (vm *VerificationManager) NewTokenForPhone() (string, error) {
	if vm.isTesting {
		return strings.Repeat("1", domain.PhoneTokenLen), nil
	}

	return domain.GenerateRandomNumbers(domain.PhoneTokenLen)
}

// DataForNewVerificationOTP генерирует и нормализует данные для последующей верификации смс-токена
func (vm *VerificationManager) DataForNewVerificationOTP(
	userData UserDataForTokenGeneration,
	shouldCheckUserExistenceByEmail bool,
	penaltyTimeWithoutTDID time.Duration,
	phoneSearchFn func(phone string) (*models.User, error),
) (VerificationTokenData, error) {
	var data VerificationTokenData

	normPhone := utils.NormPhoneNum(userData.Phone)
	if normPhone == "" {
		return data, ErrPhoneNotSpecified
	}
	data.NormPhone = normPhone

	if userData.TDID != "" {
		id, err := primitive.ObjectIDFromHex(userData.TDID)
		if err != nil {
			return data, err
		}

		// находим пользователя
		user, err := vm.users.ByTDID(id)
		if !shouldCheckUserExistenceByEmail && err != nil {
			return data, err
		}

		if shouldCheckUserExistenceByEmail && err != nil {
			user, err = vm.users.ByEmail(userData.Email)
			if err != nil {
				return data, err
			}
		}
		data.User = user

		// данный запрос на восстановление пришел с tdid, перестраховываемся
		// выставляя повторную попытку вызова через 20 секунд
		data.NextAttemptAt = time.Now().In(time.UTC).Add(vm.spamPenalty)
	} else {
		user, err := phoneSearchFn(normPhone)
		if !shouldCheckUserExistenceByEmail && err != nil {
			return data, err
		}

		if shouldCheckUserExistenceByEmail && err != nil {
			user, err = vm.users.ByEmail(userData.Email)
			if err != nil {
				return data, err
			}
		}
		data.User = user

		// данный запрос на восстановление пришел без tdid, перестраховываемся
		// выставляя повторную попытку вызова через заданное время
		data.NextAttemptAt = time.Now().In(time.UTC).Add(penaltyTimeWithoutTDID)
	}

	return data, nil
}
