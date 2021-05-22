package v1

import (
	"encoding/json"
	"net/http"

	"github.com/JetBrainer/sso/internal/ports/http/resources"
	"github.com/go-chi/render"
)

type ResetVerifyRequest struct {
	Phone string `json:"phone" validate:"required,is_phone"`
}

// @Summary Очищение объекта верификации
// @Description Снимает штрафы по верификации
// @Accept json
// @Produce json
// @Tags verify
// @Param body body ResetVerifyRequest true "Данные для снятия штрафов"
// @Success 200
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 422 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/verify/reset [put]
func (a AuthResource) ResetVerify(permissionsToCheck []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		req := new(ResetVerifyRequest)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			_ = render.Render(w, r, resources.BadRequest(err))
			return
		}

		if err := a.validate.Struct(req); err != nil {
			_ = render.Render(w, r, resources.UnprocessableEntity(err))
			return
		}

		if err := a.authManager.VerificationManager().ResetVerify(r.Context(), req.Phone); err != nil {
			_ = render.Render(w, r, resources.Internal(err))
			return
		}
	}
}
