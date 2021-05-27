package auth

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SignInByLogin осуществляет логику входа пользователя в систему по
// имени пользователя. SignInByLogin() должен уметь обрабатывать ввод
// неверных паролей, неверных имен пользователей, неверные права,
// выключенного пользователя. Функция возвращает глобальный идентификатор
// пользователя и ошибку.
func (lm *LoginManager) SignInByLogin(login, password string) (string, error) {
	// совпадает ли хэш? существует ли такой пользователь?
	valid, _ := lm.users.CheckPasswordByLogin(login, password)
	// Не нужно обрабатывать error от users.CheckPasswordByTDID(), дабы при
	// неверном имени пользователя иметь ту же самую ошибку ErrInvalidLoginOrPassword, что
	// и при неправильном пароле. Тем самым мы избавляемся от проблемы брутфорса логинов.

	// если пользователь не валидный, посылаем
	if !valid {
		return "", ErrInvalidLoginOrPassword
	}

	user, err := lm.users.ByLogin(login)
	if err != nil {
		return "", err
	}

	// проверяем, если пользователь включен
	if !user.Enabled {
		return user.ID.Hex(), ErrUserDisabled
	}

	return user.ID.Hex(), nil
}

func (lm *LoginManager) SignInByTDID(tdid primitive.ObjectID, password string) (string, error) {
	// совпадает ли хэш? существует ли такой пользователь?
	valid, _ := lm.users.CheckPasswordByTDID(tdid, password)
	// Не нужно обрабатывать error от users.CheckPasswordByTDID(), дабы при
	// неверном имени пользователя иметь ту же самую ошибку ErrInvalidLoginOrPassword, что
	// и при неправильном пароле. Тем самым мы избавляемся от проблемы брутфорса логинов.

	// если пользователь не валидный, посылаем
	if !valid {
		return "", ErrInvalidLoginOrPassword
	}

	user, err := lm.users.ByTDID(tdid)
	if err != nil {
		return "", err
	}

	// проверяем, если пользователь включен
	if !user.Enabled {
		return user.ID.Hex(), ErrUserDisabled
	}

	return user.ID.Hex(), nil
}

// SignInByEmail осуществляет логику входа пользователя в систему по его email.
func (lm *LoginManager) SignInByEmail(email, password string) (string, error) {
	user, err := lm.users.ByEmail(email)
	if err != nil {
		return "", err
	}

	return lm.SignInByTDID(user.ID, password)
}
