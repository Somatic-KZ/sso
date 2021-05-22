package monitoring

import (
	"fmt"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
)

type Health struct {
	db drivers.DataStore
}

// CheckDatabase проверка работоспособности БД
func (h *Health) CheckDatabase() error {
	if h.db == nil {
		return fmt.Errorf("[ERROR] Database instance not created")
	}

	return h.db.Ping()
}

// CheckNotificator проверка работоспособности нотификатора
func (h *Health) CheckNotificator() error {
	return nil
}

// CheckRestoreDirector проверка работоспособности директора
func (h *Health) CheckRestoreDirector() error {
	return nil
}
