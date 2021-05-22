package database

import (
	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"github.com/JetBrainer/sso/internal/adapters/database/drivers/mongo"
)

func New(conf drivers.DataStoreConfig) (drivers.DataStore, error) {
	if conf.DataStoreName == "mongo" {
		return mongo.New(conf)
	}

	return nil, ErrDatastoreNotImplemented
}
