package v1

import (
	"net/http"

	"github.com/JetBrainer/sso/internal/domain/models"
	"github.com/JetBrainer/sso/internal/domain/models/api"
	"github.com/JetBrainer/sso/internal/ports/http/resources"
	"github.com/JetBrainer/sso/internal/ports/http/resources/profile"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Профиль
// @Description Позволяет пользователю получить информацию о своем аккаунте со всеми правами и ролями
// @Produce json
// @Tags profile
// @Security JWT
// @Success 200 {object} api.ProfileResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /profile [get]
func (p ProfileResource) Profile(w http.ResponseWriter, r *http.Request) {
	var tdid interface{}
	if tdid = r.Context().Value("tdid"); tdid == "" {
		_ = render.Render(w, r, resources.BadRequest(profile.ErrUnknownTDID))
		return
	}

	id, err := primitive.ObjectIDFromHex(tdid.(string))
	if err != nil {
		_ = render.Render(w, r, resources.BadRequest(err))
		return
	}

	users := p.authManager.Users()

	var user *models.User

	user, err = users.ByTDID(id)
	if err != nil {
		_ = render.Render(w, r, resources.ResourceNotFound(err))
		return
	}

	strID, _ := tdid.(string)
	render.JSON(w, r, api.ProfileResponse{
		Created:      user.Created,
		Updated:      user.Updated,
		BirthDate:    user.BirthDate,
		Roles:        user.Roles,
		Receivers:    user.Receivers,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Patronymic:   user.Patronymic,
		Email:        user.Email,
		PrimaryPhone: user.PrimaryPhone,
		Phones:       user.Phones,
		Language:     user.Language,
		IIN:          user.IIN,
		Sex:          user.Sex,
		ID:           strID,
	})
}
