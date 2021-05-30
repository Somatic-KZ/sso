package v1

import (
	"encoding/json"
	"net/http"

	"github.com/JetBrainer/sso/internal/domain/models/api"
	"github.com/JetBrainer/sso/internal/ports/http/resources"
	"github.com/JetBrainer/sso/internal/ports/http/resources/auth"
	"github.com/JetBrainer/sso/utils"
	"github.com/go-chi/render"
)

// @Summary Регистрация
// @Description Создает нового пользователя в SSO
// @Accept json
// @Produce json
// @Tags auth
// @Param body body api.SignUPRequest true "Необходимые данные для регистрации"
// @Success 200 {object} api.SignUPResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 422 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/signup [put]
func (a AuthResource) SignUP(w http.ResponseWriter, r *http.Request) {
	var sur api.SignUPRequest

	// метрики: общее количество запросов
	if a.metrics != nil && a.metrics.SignUpByPhoneRequests != nil {
		(*a.metrics.SignUpByPhoneRequests).Add(1)
	}

	if err := json.NewDecoder(r.Body).Decode(&sur); err != nil {
		_ = render.Render(w, r, resources.BadRequest(err))
		return
	}

	//if err := a.validate.Struct(sur); err != nil {
	//	if a.metrics != nil && a.metrics.SignUpValidationErrors != nil {
	//		(*a.metrics.SignUpValidationErrors).Add(1)
	//	}
	//
	////	uniqueRules := []string{"unique_email", "unique_phones", "unique_phone"}
	////	for _, rule := range uniqueRules {
	////		if strings.Contains(err.Error(), rule) {
	////			http.Error(w, errors.New("validation error").Error(), http.StatusBadRequest )
	////			return
	////		}
	////	}
	//
	//	http.Error(w,"unprocessable value", http.StatusUnprocessableEntity)
	//	return
	//}

	normPhone := utils.NormPhoneNum(sur.Phone)
	if normPhone == "" {
		_ = render.Render(w, r, resources.UnprocessableEntity(auth.ErrInvalidNumberFormat))
		return
	}

	user, err := a.authManager.Users().CreateFromSUR(&sur)
	if err != nil {
		if a.metrics != nil && a.metrics.SignUpInternalServerErrors != nil {
			(*a.metrics.SignUpInternalServerErrors).Add(1)
		}
		_ = render.Render(w, r, resources.Internal(auth.ErrServerProblem))
		return
	}

	render.Status(r, http.StatusCreated)
	if a.metrics != nil && a.metrics.SignUpSuccessfulOperation != nil {
		(*a.metrics.SignUpSuccessfulOperation).Add(1)
	}

	// все ок, возвращаем созданный ID пользователя и ждем верификации по этому TDID
	render.JSON(w, r, api.SignUPResponse{
		TDID:   user.ID.Hex(),
		Status: "SignUP successful. Verification Required",
	})
}
