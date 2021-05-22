package models

import (
	"time"

	"github.com/JetBrainer/sso/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const defaultLanguage = "ru"

// User модель данных пользователя.
type User struct {
	Verify       Verify             `bson:"verify" json:"verify"`
	Created      time.Time          `bson:"created" json:"created"`
	Updated      time.Time          `bson:"updated" json:"updated"`
	Roles        []string           `bson:"roles" json:"roles"`
	Receivers    []Receiver         `bson:"receivers" json:"receivers,omitempty"`
	Phones       []string           `bson:"phones" json:"phones"`
	PrimaryPhone string             `bson:"phone" json:"phone"`
	Login        string             `bson:"login" json:"login"`
	FirstName    string             `bson:"firstname,omitempty" json:"firstname,omitempty"`
	LastName     string             `bson:"lastname,omitempty" json:"lastname,omitempty"`
	Patronymic   string             `bson:"patronymic,omitempty" json:"patronymic,omitempty"`
	Password     string             `bson:"password" json:"password"`
	Email        string             `bson:"email" json:"email"`
	Website      string             `bson:"website,omitempty" json:"website,omitempty"`
	Language     string             `bson:"lang,omitempty" json:"lang,omitempty"`
	Sex          string             `bson:"sex,omitempty" json:"sex,omitempty"`
	BirthDate    *time.Time         `bson:"birth_date,omitempty" json:"birth_date,omitempty"`
	Restore      *Restore           `bson:"restore,omitempty" json:"restore,omitempty"`
	IIN          int                `bson:"iin,omitempty" json:"iin,omitempty"`
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Enabled      bool               `bson:"enabled" json:"enabled"`
	Devices      []Device           `bson:"devices,omitempty" json:"devices,omitempty"`
	BankData     *BankData          `bson:"bankData,omitempty" json:"bankData,omitempty"`
	LastVerified *LastVerified      `bson:"last_verified,omitempty" json:"last_verified,omitempty"`
}

type UserShortInfo struct {
	ID         primitive.ObjectID `json:"id" validate:"required"`
	Phone      string             `json:"phone"`
	Phones     []string           `json:"phones"`
	Email      string             `json:"email"`
	Firstname  string             `json:"firstname"`
	Lastname   string             `json:"lastname"`
	Patronymic string             `json:"patronymic,omitempty"`
	BirthDate  *time.Time         `json:"birthDate"`
	IIN        int                `json:"iin,omitempty"`
	Sex        string             `json:"sex,omitempty"`
	Roles      []string           `json:"roles"`
	Receivers  []Receiver         `json:"receivers"`
	Created    time.Time          `json:"created"`
	Updated    time.Time          `json:"updated"`
	Enabled    bool               `json:"enabled"`
}

type LastVerified struct {
	Phone  string  `bson:"phone" json:"phone"`
	Email  string  `bson:"email" json:"email"`
	Device *Device `bson:"device" json:"device"`
}

// NewGeneralUser создает физического пользователя с минимально
// необходимой информацией.
func NewGeneralUser(email, phone, password, firstName, lastName, patronymic, lang string) *User {
	if lang == "" {
		lang = defaultLanguage
	}

	return &User{
		ID:           primitive.NewObjectID(),
		Login:        email,
		Password:     password,
		FirstName:    firstName,
		LastName:     lastName,
		Patronymic:   patronymic,
		Email:        email,
		PrimaryPhone: utils.NormPhoneNum(phone),
		Phones:       make([]string, 0),
		Language:     lang,
		Verify:       Verify{},
		Enabled:      true,
		Created:      time.Now().In(time.UTC),
		Updated:      time.Now().In(time.UTC),
		Roles:        []string{"user"},
		Receivers:    make([]Receiver, 0),
		Devices:      make([]Device, 0),
	}
}

// normPhoneNum нормализует номер телефона.
func (u *User) normPhoneNum(phone string) string {
	return utils.NormPhoneNum(phone)
}

// PhoneInPhones возвращает true, если запрошенный номер принадлежит пользователю.
func (u *User) PhoneInPhones(phone string) bool {
	normPhone := u.normPhoneNum(phone)
	for _, p := range u.Phones {
		if normPhone == p {
			return true
		}
	}

	return false
}

func (u User) GetShortInfo() *UserShortInfo {
	return &UserShortInfo{
		ID:         u.ID,
		Phone:      u.PrimaryPhone,
		Phones:     u.Phones,
		Email:      u.Email,
		Firstname:  u.FirstName,
		Lastname:   u.LastName,
		Patronymic: u.Patronymic,
		BirthDate:  u.BirthDate,
		IIN:        u.IIN,
		Sex:        u.Sex,
		Roles:      u.Roles,
		Receivers:  u.Receivers,
		Created:    u.Created,
		Updated:    u.Updated,
		Enabled:    u.Enabled,
	}
}

func (u *User) VerifiedContain(phone string) bool {
	for _, p := range u.Phones {
		if phone == p {
			return true
		}
	}
	return false
}

type MagentoUser struct {
	Firstname string `db:"firstname"`
	Lastname  string `db:"lastname"`
	Phone     string `db:"phone"`
}
