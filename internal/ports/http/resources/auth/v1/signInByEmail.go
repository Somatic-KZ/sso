package v1

import (
	"encoding/json"
	"net/http"

	"github.com/JetBrainer/sso/internal/domain/models/api"
	"github.com/JetBrainer/sso/internal/ports/http/resources"
	"github.com/go-chi/render"
)

// @Summary Вход по email
// @Description Проверяет пользовательский email и пароль, выписывает JWT токен
// @Accept json
// @Produce json
// @Tags auth
// @Param body body api.SignInByEmailRequest true "Необходимые данные для аутентификации"
// @Success 200 {object} api.NewJWTTokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/signin/email [post]
func (a AuthResource) SignInByEmail(w http.ResponseWriter, r *http.Request) {
	var creds api.SignInByEmailRequest

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		_ = render.Render(w, r, resources.BadRequest(err))
		return
	}

	loginManager := a.authManager.LoginManager()

	// проверяем наличие пользователя и его пароль
	tdid, err := loginManager.SignInByEmail(creds.Email, creds.Password)
	if err != nil {
		_ = render.Render(w, r, resources.Unauthorized(err))
		return
	}

	token, err := a.authManager.NewAccessToken(tdid)
	if err != nil {
		_ = render.Render(w, r, resources.Internal(err))
		return
	}

	refreshToken, err := a.authManager.NewRefreshToken(tdid)
	if err != nil {
		_ = render.Render(w, r, resources.Internal(err))
		return
	}

	render.JSON(w, r, api.NewJWTTokenResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		Status:       "SignIn success",
	})
}
