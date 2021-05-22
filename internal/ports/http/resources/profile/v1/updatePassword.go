package v1

import (
	"encoding/json"
	"net/http"

	"github.com/JetBrainer/sso/internal/domain/models/api"
	"github.com/JetBrainer/sso/internal/ports/http/resources"
	"github.com/JetBrainer/sso/internal/ports/http/resources/profile"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Смена пароля
// @Description Позволяет обновить пользователю пароль
// @Accept json
// @Produce json
// @Tags profile
// @Security JWT
// @Param body body api.UpdatePasswordRequest true "Новый пароль"
// @Success 200
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 422 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /profile/password [put]
func (p ProfileResource) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var tdid interface{}
	if tdid = r.Context().Value("tdid"); tdid == "" {
		_ = render.Render(w, r, resources.BadRequest(profile.ErrUnknownTDID))
		return
	}

	var request api.UpdatePasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		_ = render.Render(w, r, resources.BadRequest(err))
		return
	}

	if err := p.validate.Struct(request); err != nil {
		_ = render.Render(w, r, resources.UnprocessableEntity(err))
		return
	}

	id, err := primitive.ObjectIDFromHex(tdid.(string))
	if err != nil {
		_ = render.Render(w, r, resources.BadRequest(err))
		return
	}

	users := p.authManager.Users()

	if err = users.UpdatePassword(id, request.Password); err != nil {
		_ = render.Render(w, r, resources.ResourceNotFound(err))
		return
	}


	render.Status(r, http.StatusOK)
}
