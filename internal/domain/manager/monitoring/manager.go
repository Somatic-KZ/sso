package monitoring

import (
	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"github.com/JetBrainer/sso/internal/domain/models"
)

type Manager struct {
	db      drivers.DataStore
	metrics *models.Metrics
}

func New(db drivers.DataStore, metrics *models.Metrics) *Manager {
	return &Manager{
		db:      db,
		metrics: metrics,
	}
}

func (m *Manager) Health() *Health {
	return &Health{db: m.db}
}

func (m *Manager) Metrics() *models.Metrics {
	return m.metrics
}
