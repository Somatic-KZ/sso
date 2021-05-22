package auth

import (
	"log"
	"time"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"github.com/JetBrainer/sso/internal/domain/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

type Manager struct {
	db               drivers.DataStore
	jwtKey           []byte
	TokenTTL         time.Duration
	RefreshTokenTTL  time.Duration
	DelegateTokenTTL time.Duration
	isTesting        bool
}

// Claims структура, хранящая закодированный JWT авторизации.
// Встраивание для предоставления поля expiry.
type Claims struct {
	TDID      string   `json:"tdid"`
	IsRefresh bool     `json:"is_refresh"`
	Roles     []string `json:"roles"`
	jwt.StandardClaims
}

func New(db drivers.DataStore, jwtKey []byte, tokenTTL, refreshTokenTTL, delegateTokenTTL time.Duration) *Manager {
	return &Manager{
		db:               db,
		jwtKey:           jwtKey,
		TokenTTL:         tokenTTL,
		RefreshTokenTTL:  refreshTokenTTL,
		DelegateTokenTTL: delegateTokenTTL,
	}
}

func (m *Manager) Testing() {
	m.isTesting = true
}

// JWT key
func (m *Manager) JWTKey() []byte {
	return m.jwtKey
}

func (m *Manager) TokenAuth() *jwtauth.JWTAuth {
	return jwtauth.New("HS256", m.jwtKey, nil)
}

func (m Manager) NewRefreshToken(tdid string) (string, error) {
	claims := new(Claims)
	claims.TDID = tdid
	claims.IsRefresh = true
	claims.ExpiresAt = time.Now().Add(m.RefreshTokenTTL).Unix()

	userTDID, err := primitive.ObjectIDFromHex(tdid)
	if err != nil {
		log.Printf("Could not convert %s into primitive object id", tdid)
		return "", err
	}
	user, err := m.Users().ByTDID(userTDID)
	if err != nil {
		log.Printf("Could not find user with tdid %s to add roles into JWT token", tdid)
		return "", err
	}
	claims.Roles = user.Roles

	_, tokenString, err := m.TokenAuth().Encode(claims)

	return tokenString, err
}

func (m Manager) NewAccessToken(tdid string) (string, error) {
	claims := new(Claims)
	claims.TDID = tdid
	claims.ExpiresAt = time.Now().Add(m.TokenTTL).Unix()

	userTDID, err := primitive.ObjectIDFromHex(tdid)
	if err != nil {
		log.Printf("Could not convert %s into primitive object id", tdid)
		return "", err
	}
	user, err := m.Users().ByTDID(userTDID)
	if err != nil {
		log.Printf("Could not find user with tdid %s to add roles into JWT token", tdid)
		return "", err
	}
	claims.Roles = user.Roles

	_, tokenString, err := m.TokenAuth().Encode(claims)

	return tokenString, err
}


func userToClaims(u *models.User) UserClaims {
	return UserClaims{
		ID:           u.ID,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Email:        u.Email,
		PrimaryPhone: u.PrimaryPhone,
		Phones:       u.Phones,
	}
}

// Users осуществляет примитивы для работы с пользователями.
type Users struct {
	db drivers.DataStore
}

// Users создает менеджер по управлению пользователями.
func (m *Manager) Users() *Users {
	return &Users{db: m.db}
}

type LoginManager struct {
	db    drivers.DataStore
	users *Users
}

func (m *Manager) LoginManager() *LoginManager {
	return &LoginManager{
		db:    m.db,
		users: m.Users(),
	}
}
