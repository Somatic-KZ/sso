package auth

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"github.com/JetBrainer/sso/internal/domain"
	errors2 "github.com/JetBrainer/sso/internal/domain/errors"
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

// WithLongSpamPenalty устанавливает время пенальти при попытках спама
// восстановления токена без указанного tdid.
func (vm *VerificationManager) WithLongSpamPenalty(t time.Duration) *VerificationManager {
	vm.longSpamPenalty = t
	return vm
}

// AddNewPhoneAndNewVerificationToken добавляет пользоваетлю новый номер телефона и
// запускает его процесс верификации.
func (vm *VerificationManager) AddNewPhoneAndNewVerificationToken(ctx context.Context, tdid, phone string) error {
	normPhone := utils.NormPhoneNum(phone)
	if normPhone == "" {
		return ErrPhoneNotSpecified
	}

	if tdid == "" {
		return drivers.ErrUserIDNotSpec
	}

	id, err := primitive.ObjectIDFromHex(tdid)
	if err != nil {
		return err
	}

	// находим пользователя
	user, err := vm.users.ByTDID(id)
	if err != nil {
		return err
	}

	// проверяем, может быть номер уже был ранее верифицирован?
	if user.PhoneInPhones(normPhone) {
		return errors.New("Phone already verified")
	}

	// может быть пользователь многократно вызывает сервис, нельзя допустить
	// многократную отсылку SMS, даем на запрос 10 секунд.
	if user.Verify.Expired.Add(-VerifyTTL).In(time.UTC).After(time.Now().Add(-vm.spamPenalty).In(time.UTC)) {
		return errors2.ErrTooManyRequests
	}

	user.PrimaryPhone = normPhone
	if err := vm.db.UserUpdate(ctx, user); err != nil {
		return err
	}

	var token string

	token, err = vm.NewTokenForPhone()
	if err != nil {
		return err
	}

	return vm.db.VerifyPhone(
		ctx,
		user.ID,
		normPhone,
		token,
		time.Now().In(time.UTC).Add(VerifyTTL),
		time.Now().In(time.UTC).Add(vm.spamPenalty),
	)
}

// NewVerificationToken запускает процесс верификации номера телефона.
func (vm *VerificationManager) NewVerificationToken(ctx context.Context, tdid, phone string) error {
	userData := UserDataForTokenGeneration{
		TDID:  tdid,
		Phone: phone,
	}
	data, err := vm.DataForNewVerificationOTP(userData, false, vm.longSpamPenalty, vm.users.ByPrimaryPhone)
	if err != nil {
		return err
	}
	user, normPhone, nextAttemptAt := data.User, data.NormPhone, data.NextAttemptAt

	// проверяем, может быть номер уже был ранее верифицирован?
	if user.PhoneInPhones(normPhone) {
		return errors.New("Phone already verified")
	}

	// валидируем принадлежность номера телефона пользователю
	if user.PrimaryPhone != normPhone {
		return ErrPhoneNumNotLinkedToAccount
	}

	// может быть пользователь многократно вызывает сервис, нельзя допустить
	// многократную отсылку SMS
	if user.Verify.NextAttemptAt.In(time.UTC).After(time.Now().In(time.UTC)) {
		return errors2.ErrTooManyRequests
	}

	token, err := vm.NewTokenForPhone()
	if err != nil {
		return err
	}
	return vm.db.VerifyPhone(
		ctx,
		user.ID,
		normPhone,
		token,
		time.Now().In(time.UTC).Add(VerifyTTL),
		nextAttemptAt,
	)
}

// NewDelegateToken запускает процесс верификации номера телефона.
func (vm *VerificationManager) NewDelegateToken(ctx context.Context, phone string) error {
	userData := UserDataForTokenGeneration{
		Phone: phone,
	}
	data, err := vm.DataForNewVerificationOTP(
		userData,
		false,
		vm.spamPenalty,
		vm.users.ByAnyPhone,
	)
	if err != nil {
		return err
	}
	user, normPhone := data.User, data.NormPhone

	// может быть пользователь многократно вызывает сервис, нельзя допустить
	// многократную отсылку SMS
	if user.Verify.NextAttemptAt.In(time.UTC).After(time.Now().In(time.UTC)) {
		return errors2.ErrTooManyRequests
	}

	token, err := vm.NewTokenForPhone()
	if err != nil {
		return err
	}

	return vm.db.VerifyPhoneIncrementGeneration(
		ctx,
		user.ID,
		normPhone,
		token,
		VerifyTTL,
		vm.spamPenalty,
	)
}

// NewFastSignInToken запускает процесс верификации номера телефона
// для последующей "быстрой" авторизации.
func (vm *VerificationManager) NewFastSignInToken(ctx context.Context, tdid, phone, email string) (string, error) {
	userData := UserDataForTokenGeneration{
		TDID:  tdid,
		Phone: phone,
		Email: email,
	}
	data, err := vm.DataForNewVerificationOTP(userData, true, vm.spamPenalty*3, vm.users.ByAnyPhone)
	if err != nil {
		return "", err
	}
	user, normPhone, nextAttemptAt := data.User, data.User.PrimaryPhone, data.NextAttemptAt

	// может быть пользователь многократно вызывает сервис, нельзя допустить
	// многократную отсылку SMS
	if user.Verify.NextAttemptAt.In(time.UTC).After(time.Now().In(time.UTC)) {
		return "", errors2.ErrTooManyRequests
	}

	token, err := vm.NewTokenForPhone()
	if err != nil {
		return "", err
	}

	err = vm.db.VerifyPhone(
		ctx,
		user.ID,
		normPhone,
		token,
		time.Now().In(time.UTC).Add(VerifyTTL),
		nextAttemptAt,
	)
	if err != nil {
		return "", err
	}

	return utils.MaskPhoneNum(normPhone), nil
}

// NewTokenForPhone генерирует новый токен.
func (vm *VerificationManager) NewTokenForPhone() (string, error) {
	if vm.isTesting {
		return strings.Repeat("1", domain.PhoneTokenLen), nil
	}

	return domain.GenerateRandomNumbers(domain.PhoneTokenLen)
}

func (vm *VerificationManager) ValidateToken(ctx context.Context, tdid, phone, token string) (*models.User, error) {
	if tdid == "" && phone == "" {
		return nil, drivers.ErrUserIDNotSpec
	}

	if token == "" {
		return nil, drivers.ErrTokenNotSpec
	}

	var user *models.User
	var err error
	var id primitive.ObjectID

	if tdid != "" {
		id, err = primitive.ObjectIDFromHex(tdid)
		if err != nil {
			return nil, err
		}

		// находим пользователя
		user, err = vm.users.ByTDID(id)
		if err != nil {
			return nil, err
		}
	}

	if phone != "" {
		user, err = vm.users.ByPrimaryPhone(utils.NormPhoneNum(phone))
		if err != nil {
			return nil, err
		}

		id = user.ID
	}

	if user.Verify.Token != token {
		if user.Verify.Tries == TriesLimit {
			// удаляем токен на восстановления
			if err := vm.db.VerifyClean(context.Background(), user.ID, phone); err != nil {
				log.Printf("[ERROR] VerifyClean() failed: %v", err)
			}

			return nil, errors2.ErrTokenDoesNotExist
		}

		if user.Verify.Token != "" {
			if err := vm.db.VerifyIncrementTries(context.Background(), user.ID); err != nil {
				log.Printf("[ERROR] VerifyIncrementTries() failed: %v", err)
			}
		}

		return nil, errors2.ErrTokenDoesNotExist
	}

	if user.Verify.Phone == "" {
		return nil, ErrPhoneNotSpecified
	}

	// проверяем, не заэкспайрился ли токен
	if user.Verify.Expired.In(time.UTC).Before(time.Now().In(time.UTC)) {
		return nil, errors2.ErrTokenHasExpired
	}

	if err := vm.db.VerifyClean(ctx, id, user.Verify.Phone); err != nil {
		return nil, err
	}

	return user, nil
}

func (vm *VerificationManager) ValidateFastSignInToken(
	ctx context.Context,
	tdid,
	phone,
	email,
	token string,
) (*models.User, error) {
	if tdid == "" && phone == "" {
		return nil, drivers.ErrUserIDNotSpec
	}

	if token == "" {
		return nil, drivers.ErrTokenNotSpec
	}

	var user *models.User
	var err error
	var id primitive.ObjectID

	if tdid != "" {
		id, err = primitive.ObjectIDFromHex(tdid)
		if err != nil {
			return nil, err
		}

		// находим пользователя
		user, err = vm.users.ByTDID(id)

		if email == "" && err != nil {
			return nil, err
		}

		// Если юзер предоставил имейл и мы его не нашли по TDID, пытаемся найти по имейлу
		if email != "" && err != nil {
			user, err = vm.users.ByEmail(email)
			if err != nil {
				return nil, err
			}
		}
	}

	if phone != "" && tdid == "" {
		user, err = vm.users.ByAnyPhone(utils.NormPhoneNum(phone))
		if err != nil && email == "" {
			return nil, err
		}

		// Если юзер предоставил имейл и мы его не нашли по телефону, пытаемся найти по имейлу
		if email != "" && err != nil {
			user, err = vm.users.ByEmail(email)
			if err != nil {
				return nil, err
			}
		}

		id = user.ID
	}

	if user.Verify.Token != token {
		if user.Verify.Tries == TriesLimit {
			// удаляем токен на восстановления
			if err := vm.db.VerifyClean(context.Background(), user.ID, user.PrimaryPhone); err != nil {
				log.Printf("[ERROR] VerifyClean() failed: %v", err)
			}

			return nil, errors2.ErrTokenDoesNotExist
		}

		if user.Verify.Token != "" {
			if err := vm.db.VerifyIncrementTries(context.Background(), user.ID); err != nil {
				log.Printf("[ERROR] VerifyIncrementTries() failed: %v", err)
			}
		}

		return nil, errors2.ErrTokenDoesNotExist
	}

	if user.Verify.Phone == "" {
		return nil, ErrPhoneNotSpecified
	}

	// проверяем, не заэкспайрился ли токен
	if user.Verify.Expired.In(time.UTC).Before(time.Now().In(time.UTC)) {
		return nil, errors2.ErrTokenHasExpired
	}

	if err := vm.db.VerifyClean(ctx, id, user.Verify.Phone); err != nil {
		return nil, err
	}

	return user, nil
}

func (vm *VerificationManager) ValidatePhoneSignIn(ctx context.Context, phone, otp string) (*models.User, error) {
	if phone == "" {
		return nil, errors2.ErrUserIDNotSpec
	}

	if otp == "" {
		return nil, errors2.ErrOTPNotSpec
	}

	var user *models.User
	var err error
	var id primitive.ObjectID

	user, err = vm.users.ByAnyPhone(utils.NormPhoneNum(phone))
	if err != nil {
		return nil, err
	}

	id = user.ID

	if user.Verify.Token != otp {
		if user.Verify.Tries == TriesLimit {
			// удаляем токен на восстановления
			if err := vm.db.VerifyClean(context.Background(), user.ID, user.PrimaryPhone); err != nil {
				log.Printf("[ERROR] VerifyClean() failed: %v", err)
			}

			return nil, errors2.ErrTokenDoesNotExist
		}

		if user.Verify.Token != "" {
			if err := vm.db.VerifyIncrementTries(context.Background(), user.ID); err != nil {
				log.Printf("[ERROR] VerifyIncrementTries() failed: %v", err)
			}
		}

		return nil, errors2.ErrTokenDoesNotExist
	}

	if user.Verify.Phone == "" {
		return nil, ErrPhoneNotSpecified
	}

	// проверяем, не заэкспайрился ли токен
	if user.Verify.Expired.In(time.UTC).Before(time.Now().In(time.UTC)) {
		return nil, errors2.ErrTokenHasExpired
	}

	if err := vm.db.VerifyClean(ctx, id, user.Verify.Phone); err != nil {
		return nil, err
	}

	return user, nil
}

func (vm *VerificationManager) FindNewVerifyRequests(ctx context.Context, c chan<- models.User) {
	vm.db.VerifyFindNew(ctx, c)
}

func (vm *VerificationManager) VerifySendNotificationSuccessfully(ctx context.Context, tdid primitive.ObjectID) error {
	return vm.db.VerifySendNotificationSuccessfully(ctx, tdid)
}

func (vm VerificationManager) ResetVerify(ctx context.Context, phone string) error {
	user, err := vm.users.ByAnyPhone(utils.NormPhoneNum(phone))
	if err != nil {
		return err
	}

	user.Verify = models.Verify{
		Status:        domain.TokenStatusFinish,
		Expired:       time.Now().Add(MinusDay),
		NextAttemptAt: time.Now().Add(MinusDay),
	}

	if err := vm.db.UpdateVerify(ctx, user); err != nil {
		return err
	}

	return nil
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
