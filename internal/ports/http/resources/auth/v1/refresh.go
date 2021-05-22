package v1

import (
	"encoding/json"
	"net/http"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	errors2 "github.com/JetBrainer/sso/internal/domain/errors"
	"github.com/JetBrainer/sso/internal/domain/manager/auth"
	"github.com/JetBrainer/sso/internal/domain/models/api"
	"github.com/JetBrainer/sso/internal/ports/http/resources"
	"github.com/JetBrainer/sso/internal/ports/http/resources/profile"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/render"
)

// @Summary Рефреш токена
// @Description Обновляет время жизни токена
// @Accept json
// @Produce json
// @Tags auth
// @Param body body api.RefreshRequest true "Необходимые данные для обновления токена"
// @Success 200 {object} api.NewJWTTokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/refresh [post]
func (a AuthResource) Refresh(w http.ResponseWriter, r *http.Request) {
	var request api.RefreshRequest

	// метрики: общее количество запросов
	if a.metrics != nil && a.metrics.RefreshRequests != nil {
		(*a.metrics.RefreshRequests).Add(1)
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		_ = render.Render(w, r, resources.BadRequest(err))
		return
	}

	if request.RefreshToken == "" {
		_ = render.Render(w, r, resources.BadRequest(drivers.ErrTokenNotSpec))
		return
	}

	// инициализируем новый инстанс `Claims`
	claims := new(auth.Claims)

	// Парсим JWT строку и сохраняем результат в `claims`.
	tkn, err := jwt.ParseWithClaims(request.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return a.authManager.JWTKey(), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			_ = render.Render(w, r, resources.Unauthorized(err))
			return
		}

		_ = render.Render(w, r, resources.BadRequest(err))
		return
	}

	if !tkn.Valid {
		if a.metrics != nil && a.metrics.RefreshInvalidTokenErrors != nil {
			(*a.metrics.RefreshInvalidTokenErrors).Add(1)
		}
		_ = render.Render(w, r, resources.Unauthorized(errors2.ErrTokenDoesNotExist))
		return
	}

	if claims.TDID == "" {
		if a.metrics != nil && a.metrics.RefreshInvalidTokenErrors != nil {
			(*a.metrics.RefreshInvalidTokenErrors).Add(1)
		}
		_ = render.Render(w, r, resources.BadRequest(profile.ErrUnknownTDID))
		return
	}

	if !claims.IsRefresh {
		if a.metrics != nil && a.metrics.RefreshInvalidTokenErrors != nil {
			(*a.metrics.RefreshInvalidTokenErrors).Add(1)
		}
		_ = render.Render(w, r, resources.Unauthorized(ErrInvalidToken))
		return
	}

	var newToken string
	newToken, err = a.authManager.NewAccessToken(claims.TDID)
	if err != nil {
		if a.metrics != nil && a.metrics.RefreshInternalServerErrors != nil {
			(*a.metrics.RefreshInternalServerErrors).Add(1)
		}
		_ = render.Render(w, r, resources.Internal(err))
		return
	}

	if a.metrics != nil && a.metrics.RefreshSuccessfulOperation != nil {
		(*a.metrics.RefreshSuccessfulOperation).Add(1)
	}

	render.JSON(w, r, api.NewJWTTokenResponse{
		AccessToken: newToken,
		Status:      "Refresh token successful",
	})
}
