package resources

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"

	"github.com/JetBrainer/sso/internal/domain/manager/monitoring"
	"github.com/JetBrainer/sso/internal/domain/models/api"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type MonitoringResource struct {
	manager *monitoring.Manager
}

func NewHealth(manager *monitoring.Manager) *MonitoringResource {
	return &MonitoringResource{
		manager: manager,
	}
}

func (mr MonitoringResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", mr.HealthCheck)

	return r
}

// HealthCheck проверка работоспособности контейнера:
// Проверку жизни БД, нотификатора, Restore директора.
func (mr *MonitoringResource) HealthCheck(w http.ResponseWriter, r *http.Request) {
	var err error
	health := mr.manager.Health()

	err = health.CheckDatabase()
	if err != nil {
		_ = render.Render(w, r, NewHealthResponse(api.HealthStatusNotOK))
		return
	}

	err = health.CheckNotificator()
	if err != nil {
		_ = render.Render(w, r, NewHealthResponse(api.HealthStatusNotOK))
		return
	}

	err = health.CheckRestoreDirector()
	if err != nil {
		_ = render.Render(w, r, NewHealthResponse(api.HealthStatusNotOK))
		return
	}

	_ = render.Render(w, r, NewHealthResponse(api.HealthStatusOK))
}

// NewHealthResponse создает новые ответ
func NewHealthResponse(status string) *api.HealthResponse {
	return &api.HealthResponse{
		ID:      generateUID(),
		Node:    "foobar",
		Name:    "Serf Health Status",
		CheckID: "serfHealth",
		Status:  status,
	}
}

// generateUID генерирует уникальный id
func generateUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)

	if err != nil {
		log.Printf("[Error] " + err.Error())
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
