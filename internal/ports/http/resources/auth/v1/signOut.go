package v1

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

const (
	tokenTTL = 60 * time.Minute
)

// @Summary Завершает сессию пользователя, удаляя JWT токен.
// @Description Успешный вызов удаляет представленный token в cookie.
// @Produce json
// @Tags auth
// @Security JWT
// @Success 200
// @Router /auth/signout [delete]
func (a AuthResource) SignOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    "",
		HttpOnly: false,
		Expires:  time.Now().In(time.UTC).Add(-tokenTTL),
	})

	render.Status(r, http.StatusOK)
}
