package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тестирует создание хэшей по паролю.
// Тест проверяеет создание двух различных хэшей по одинаковому паролю.
func TestHashPassword(t *testing.T) {
	s := "password"

	// генерируем первый хэш
	hash1, err := HashPassword(s)
	assert.NoError(t, err)
	assert.NotNil(t, hash1)

	// генерируем второй хэш
	hash2, err := HashPassword(s)
	assert.NoError(t, err)
	assert.NotNil(t, hash2)

	// оба хэша не должны быть равны
	assert.NotEqual(t, hash1, hash2)
}

// Проверка пароля по хэшу.
// Тест сравнивает два разных хэша от одного пароля, неправильный хэш с
// правильным паролем и правильный хэш с неправильным паролем.
func TestCheckPasswordHash(t *testing.T) {
	s := "password"
	hash1 := "$2a$10$sn5zKytNAC9ARobKzCXzf.5eibqi92SJTxFrJPlkTXshGwKo/Pw2i"
	hash2 := "$2a$10$4F6CFRZ80jmW9W/qaLg.uu07LGUO/jkf7TBmdjqO4PA9H/P/A/AfO"
	badHash := "$2a$10$4F6CFRZ80jmW9W/qa12.uu07LGUO/jkf7TBmdjqO4PA9H/P/A/AfO"

	// проверяем два правильных хэша
	assert.Equal(t, true, CheckPasswordHash(s, []byte(hash1)))
	assert.Equal(t, true, CheckPasswordHash(s, []byte(hash2)))

	// проверяем неправильный хэш и правильный пароль
	assert.Equal(t, false, CheckPasswordHash(s, []byte(badHash)))

	// проверяем правильный хэш и неправильный пароль
	assert.Equal(t, false, CheckPasswordHash("pass", []byte(hash1)))
}
