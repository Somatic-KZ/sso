package api

import (
	"net/http"

	"github.com/go-chi/render"
)

const HealthStatusOK = "passing"
const HealthStatusNotOK = "not_passing"

// Health - ответ на запрос о здоровье.
type HealthResponse struct {
	ID      string `json:"id"`
	Node    string `json:"node"`
	CheckID string `json:"check_id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
}

// Render выполняем интерфейс render.Renderer
func (h *HealthResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	code := http.StatusOK
	if HealthStatusOK != h.Status {
		code = http.StatusInternalServerError
	}

	render.Status(r, code)

	return nil
}
