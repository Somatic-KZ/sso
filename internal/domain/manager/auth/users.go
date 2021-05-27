package auth

import (
	"context"
	"time"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"github.com/JetBrainer/sso/internal/domain/models"
	"github.com/JetBrainer/sso/internal/domain/models/api"
	"github.com/JetBrainer/sso/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ByTDID возвращает пользователя по его TDID.
func (u *Users) ByTDID(tdid primitive.ObjectID) (*models.User, error) {
	user, err := u.db.UserByTDID(context.Background(), tdid)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserDoesNotExist
	}

	return user, nil
}

// ByLogin возвращает пользователя по его логину.
func (u *Users) ByLogin(login string) (*models.User, error) {
	user, err := u.db.UserByLogin(context.Background(), login)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserDoesNotExist
	}

	return user, nil
}

// Exists проверяет существование пользователя по его login.
func (u *Users) Exists(login string) (bool, error) {
	user, err := u.db.UserByLogin(context.Background(), login)
	if err == drivers.ErrUserDoesNotExist {
		return false, nil
	}

	return user != nil, err
}

func (u *Users) Create(user *models.User) error {
	var tdid primitive.ObjectID

	tdid, err := u.db.UserCreate(context.Background(), user)
	if err != nil {
		return err
	}

	user.ID = tdid

	return nil
}

// CreateFromSUR создает нового пользователя проксируя данные из api.SignUPRequest.
func (u *Users) CreateFromSUR(sur *api.SignUPRequest) (*models.User, error) {
	newUser := models.NewGeneralUser(
		sur.Email,
		sur.Phone,
		sur.Password,
		sur.FirstName,
		sur.LastName,
		sur.Patronymic,
		sur.Lang,
	)

	if newUser.Password != "" {
		hash, err := HashPassword(newUser.Password)

		if err != nil {
			return nil, err
		}

		newUser.Password = string(hash)
	}

	newUser.BirthDate = appendBirthDate(sur.BirthDate)

	if sur.IIN != nil {
		newUser.IIN = *sur.IIN
	}

	if sur.Sex != nil {
		newUser.Sex = *sur.Sex
	}

	var tdid primitive.ObjectID

	tdid, err := u.db.UserCreate(context.Background(), newUser)
	if err != nil {
		return nil, err
	}

	newUser.ID = tdid

	return newUser, nil
}

// appendBirthDate
func appendBirthDate(birthDate *string) *time.Time {
	if birthDate == nil {
		return nil
	}

	bd, _ := time.Parse("", *birthDate)

	return &bd
}

// appendAddress проверка адреса
func appendAddress(sur *api.SignUPRequest) *models.Address {
	if sur.Region == nil &&
		sur.City == nil &&
		sur.Street == nil &&
		sur.Corpus == nil &&
		sur.House == nil &&
		sur.Apartment == nil &&
		sur.Zipcode == nil {

		return nil
	}

	var addr models.Address

	if sur.Region != nil {
		addr.Region = *sur.Region
	}

	if sur.City != nil {
		addr.City = *sur.City
	}

	if sur.Street != nil {
		addr.Street = *sur.Street
	}

	if sur.Corpus != nil {
		addr.Corpus = *sur.Corpus
	}

	if sur.House != nil {
		addr.House = *sur.House
	}

	if sur.Apartment != nil {
		addr.Apartment = *sur.House
	}

	if sur.Zipcode != nil {
		addr.Zipcode = *sur.Zipcode
	}

	return &addr
}

// Delete удаляет пользователя.
// Метод не только должен удалять пользователя, но также, вычищать все его токены.
func (u *Users) Delete(login string) error {
	exists, err := u.Exists(login)
	if err != nil {
		return err
	}

	// такой пользователь не существует
	if !exists {
		return ErrUserDoesNotExist
	}

	// удаление пользователя
	if err := u.db.UserDelete(context.Background(), login); err != nil {
		return err
	}

	return nil
}

// Update изменяет пользователя по его ID.
// Update() использует входную структуру models.User, но не изменяет в ней пароль.
// Метод также удаляет все токены пользователя, если удалена роль администратора
// или пользователь выключен.
func (u *Users) Update(user *models.User) error {
	oldUser, err := u.ByLogin(user.Login)
	if err != nil {
		return err
	}

	if oldUser == nil {
		return ErrUserDoesNotExist
	}

	// замена хэша пароля на старый
	user.Password = oldUser.Password

	if user.Email != "" {
		// email соотетствует формату?
		if err := utils.ValidateEmail(user.Email); err != nil {
			return ErrInvalidEmailAddress
		}
	}

	// используется ли уже указанный email?
	existsUser, _ := u.ByEmail(user.Email)

	// email который хочет взять пользователь уже у кого-то есть
	if existsUser != nil && existsUser.Login != user.Login {
		return ErrEmailAlreadyTaken
	}

	// обновляем дату изменения пользователя
	user.Updated = time.Now().In(time.UTC)

	// обновление пользователя
	if err := u.db.UserUpdate(context.Background(), user); err != nil {
		return err
	}

	return nil
}

// Enable делает пользователя активным.
func (u *Users) Enable(login string) error {
	exists, err := u.Exists(login)
	if err != nil {
		return err
	}

	// такой пользователь не существует
	if !exists {
		return ErrUserDoesNotExist
	}

	user, err := u.ByLogin(login)
	if err != nil {
		return err
	}

	user.Enabled = true

	return u.Update(user)
}

// Disable делает пользователя неактивным.
func (u *Users) Disable(login string) error {
	exists, _ := u.Exists(login)
	// такой пользователь не существует
	if !exists {
		return ErrUserDoesNotExist
	}

	user, err := u.ByLogin(login)
	if err != nil {
		return err
	}

	// выключаем пользователя
	user.Enabled = false
	if err := u.Update(user); err != nil {
		return err
	}

	return err
}

// CheckPasswordByLogin сверяет пароль с хэшем.
func (u *Users) CheckPasswordByLogin(login, password string) (bool, error) {
	// получаем пользователя
	user, err := u.ByLogin(login)
	if err != nil {
		return false, err
	}

	// сравниваем хэши
	return CheckPasswordHash(password, []byte(user.Password)), nil
}

// CheckPasswordByTDID сверяет пароль с хэшем.
func (u *Users) CheckPasswordByTDID(tdid primitive.ObjectID, password string) (bool, error) {
	// получаем пользователя
	user, err := u.ByTDID(tdid)
	if err != nil {
		return false, err
	}

	// сравниваем хэши
	return CheckPasswordHash(password, []byte(user.Password)), nil
}

// UpdatePassword валидирует сложность пароля и обновляет хэш в БД.
func (u *Users) UpdatePassword(tdid primitive.ObjectID, password string) error {
	// получаем пользователя
	user, err := u.ByTDID(tdid)
	if err != nil {
		return ErrUserDoesNotExist
	}

	// валидируем сложность пароля
	/*	if err := validatePassword(password); err != nil {
		return err
	}*/

	var hash []byte

	// генерируем хэш
	hash, err = HashPassword(password)
	if err != nil {
		return err
	}

	user.Password = string(hash)

	// сохраняем пользователя используя метод UserUpdate() именно от фабрики,
	// так как u.Update() обновляет все, кроме пароля.
	return u.db.UserUpdate(context.Background(), user)
}

// ByEmail находит пользователя по его Email.
// Используется в проверках уникальности соотношения email - аккаунт,
// а также, для восстановления пароля по email адресу.
func (u *Users) ByEmail(email string) (*models.User, error) {
	return u.db.UserByEmail(context.Background(), email)
}

// FullNamesByTDID возвращает список Имя + Фамилия по списку TDID
func (u *Users) FullNamesByTDID(ctx context.Context, rawTDIDList []string) (map[string]string, error) {
	tdidList := make([]primitive.ObjectID, 0, len(rawTDIDList))

	for _, id := range rawTDIDList {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, ErrInvalidTdIDList(id)
		}
		tdidList = append(tdidList, objectID)
	}

	fullNames, err := u.db.UserFullNamesByTDID(ctx, tdidList)
	if err != nil {
		return nil, err
	}

	return fullNames, nil
}

func (u *Users) DevicesByTDID(ctx context.Context, tdid primitive.ObjectID) ([]models.Device, error) {
	devices, err := u.db.UserDevicesByTDID(ctx, tdid)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// AddReceiver добавляет получателя юзеру и возвращает строчное представление ID
func (u *Users) AddReceiver(tdid primitive.ObjectID, receiver *models.Receiver) (string, error) {
	// получаем пользователя
	user, err := u.ByTDID(tdid)
	if err != nil {
		return "", ErrUserDoesNotExist
	}

	receiver.ID = primitive.NewObjectID()

	if len(user.Receivers) >= 1 {
		if receiver.IsDefault {
			user.Receivers = u.addDefaultReceiverToUserReceivers(*receiver, user.Receivers)
		}

		if !receiver.IsDefault {
			// Просто добавляем получателя в конец списка
			user.Receivers = append(user.Receivers, *receiver)
		}
	}

	if user.Receivers == nil || len(user.Receivers) == 0 {
		receiver.IsDefault = true // Если адресов больше нет, то первый добавляемый получает статус "по умолчанию"
		user.Receivers = []models.Receiver{*receiver}
	}

	if err := u.db.UserUpdate(context.Background(), user); err != nil {
		return "", err
	}

	return receiver.ID.Hex(), nil
}

// DeleteReceiverByID удаляет получателя из списка получателей у юзера по айди получателя
func (u *Users) DeleteReceiverByID(tdid, receiverID primitive.ObjectID) error {
	// получаем пользователя
	user, err := u.ByTDID(tdid)
	if err != nil {
		return ErrUserDoesNotExist
	}

	filteredUserReceivers := make([]models.Receiver, 0, len(user.Receivers))
	for _, r := range user.Receivers {
		if r.ID != receiverID {
			filteredUserReceivers = append(filteredUserReceivers, r)
		}
	}
	user.Receivers = filteredUserReceivers

	if err := u.db.UserUpdate(context.Background(), user); err != nil {
		return err
	}

	return nil
}

func (u *Users) UpdateReceiver(tdid primitive.ObjectID, updatedReceiver *models.UpdatedReceiver) error {
	//получаем пользователя
	user, err := u.ByTDID(tdid)
	if err != nil {
		return ErrUserDoesNotExist
	}

	userReceiversCopy := make([]models.Receiver, 0, len(user.Receivers))
	for _, r := range user.Receivers {
		if r.ID == updatedReceiver.ID {
			userReceiversCopy = append(userReceiversCopy, *u.receiverSoftUpdate(&r, updatedReceiver))
			continue
		}

		// Если получатель получил статус "по умолчанию", остальные теряют этот статус
		if updatedReceiver.IsDefault != nil && *updatedReceiver.IsDefault {
			userReceiversCopy = append(userReceiversCopy, models.Receiver{
				ID:              r.ID,
				FirstName:       r.FirstName,
				LastName:        r.LastName,
				Email:           r.Email,
				PrimaryPhone:    r.PrimaryPhone,
				AdditionalPhone: r.AdditionalPhone,
				Address:         r.Address,
				IsDefault:       false,
				IsOrganization:  r.IsOrganization,
				Organization:    r.Organization,
			})
			continue
		}

		userReceiversCopy = append(userReceiversCopy, r)
	}
	user.Receivers = userReceiversCopy

	if err := u.db.UserUpdate(context.Background(), user); err != nil {
		return err
	}

	return nil
}

func (u *Users) addDefaultReceiverToUserReceivers(defaultReceiver models.Receiver, userReceivers []models.Receiver) []models.Receiver {
	userReceiversCopy := make([]models.Receiver, 0, len(userReceivers)+1)
	userReceiversCopy = append(userReceiversCopy, defaultReceiver) // Сажаем дефолтного юзера первым

	for _, r := range userReceivers {
		// Убираем статус "по умолчанию" у других получателей
		userReceiversCopy = append(userReceiversCopy, models.Receiver{
			ID:              r.ID,
			FirstName:       r.FirstName,
			LastName:        r.LastName,
			Email:           r.Email,
			PrimaryPhone:    r.PrimaryPhone,
			AdditionalPhone: r.AdditionalPhone,
			Address:         r.Address,
			IsDefault:       false,
			IsOrganization:  r.IsOrganization,
			Organization:    r.Organization,
		})
	}

	return userReceiversCopy
}

func (u *Users) receiverSoftUpdate(receiver *models.Receiver, updatedReceiver *models.UpdatedReceiver) *models.Receiver {
	// Всегда обновляем валидированные/обязательные поля и флаги
	receiver.FirstName = updatedReceiver.FirstName
	receiver.LastName = updatedReceiver.LastName
	receiver.Email = updatedReceiver.Email
	receiver.PrimaryPhone = updatedReceiver.PrimaryPhone

	// Далее сажаем опциональные данные
	if updatedReceiver.AdditionalPhone != nil {
		receiver.AdditionalPhone = *updatedReceiver.AdditionalPhone
	}
	if updatedReceiver.Address != nil {
		receiver.Address = updatedReceiver.Address
	}
	if updatedReceiver.IsDefault != nil {
		receiver.IsDefault = *updatedReceiver.IsDefault
	}
	if updatedReceiver.IsOrganization != nil {
		receiver.IsOrganization = *updatedReceiver.IsOrganization
	}
	if updatedReceiver.Organization != nil {
		receiver.Organization = updatedReceiver.Organization
	}

	// Удаляем данные по организации, если поменяли тип получателя на "физическое лицо"
	if !receiver.IsOrganization {
		receiver.Organization = nil
	}

	return receiver
}
