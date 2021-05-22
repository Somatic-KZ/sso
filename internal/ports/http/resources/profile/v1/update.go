package v1

import (
	"encoding/json"
	"net/http"

	"github.com/JetBrainer/sso/internal/domain/models"
	"github.com/JetBrainer/sso/internal/domain/models/api"
	"github.com/JetBrainer/sso/internal/ports/http/resources"
	"github.com/JetBrainer/sso/internal/ports/http/resources/profile"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Изменение профиля
// @Description Позволяет обновить данные по своему профилю
// @Accept json
// @Produce json
// @Tags profile
// @Security JWT
// @Param body body api.ProfileUpdateRequest true "Обновленные данные профиля"
// @Success 200
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 422 {object} models.ErrorResponse
// @Router /profile [put]
func (p ProfileResource) ProfileUpdate(w http.ResponseWriter, r *http.Request) {
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

	var request api.ProfileUpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		_ = render.Render(w, r, resources.BadRequest(err))
		return
	}

	// Сохраняем ошибки валидации в переменной до получения данных по юзеру
	validationErrors := p.validate.Struct(request)

	var user *models.User

	users := p.authManager.Users()
	user, err = users.ByTDID(id)
	if err != nil {
		_ = render.Render(w, r, resources.ResourceNotFound(err))
		return
	}

	if validationErrors != nil {
		validationErrorsSlice := validationErrors.(validator.ValidationErrors)
		finalValidationErrors := make(validator.ValidationErrors, 0, len(validationErrorsSlice))
		for _, valErr := range validationErrorsSlice {
			// Пропускаем ошибку уникальности имейла, так как он уже принадлежит юзеру
			if valErr.ActualTag() == "unique_email" && *request.Email == user.Email {
				continue
			}
			// Остальные ошибки записываем
			finalValidationErrors = append(finalValidationErrors, valErr)
		}

		// Если у нас остались ошибки после фильтрации, отдаем их клиенту
		if len(finalValidationErrors) > 0 {
			_ = render.Render(w, r, resources.UnprocessableEntity(finalValidationErrors))
			return
		}
	}

	// детектим изменения, обновляем структуру user и сохраняем
	if err := users.Update(userProfileRequestMerge(&request, user)); err != nil {
		_ = render.Render(w, r, resources.BadRequest(err))
		return
	}

	render.Status(r, http.StatusOK)
}
