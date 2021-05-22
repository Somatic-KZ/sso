package v1

import (
	"encoding/json"
	"net/http"

	"github.com/JetBrainer/sso/internal/ports/http/resources"
	"github.com/go-chi/render"
)

type ResetRecoveryRequest struct {
	Phone string `json:"phone" validate:"required,is_phone"`
}

// @Summary Очищение объекта восстановления
// @Description Снимает штрафы по восстановлению
// @Accept json
// @Produce json
// @Tags recovery
// @Param body body ResetRecoveryRequest true "Данные для снятия штрафов по восстановлению"
// @Success 200
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 422 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/recovery/reset [put]
func (a AuthResource) ResetRecovery(permissionsToCheck []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		req := new(ResetRecoveryRequest)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			_ = render.Render(w, r, resources.BadRequest(err))
			return
		}

		if err := a.validate.Struct(req); err != nil {
			_ = render.Render(w, r, resources.UnprocessableEntity(err))
			return
		}

		if err := a.authManager.RecoveryManager().RestoreReset(r.Context(), req.Phone); err != nil {
			_ = render.Render(w, r, resources.Internal(err))
			return
		}
	}
}
