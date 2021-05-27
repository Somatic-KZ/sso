package grpc

import (
	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"github.com/JetBrainer/sso/internal/domain/manager/auth"
	"github.com/JetBrainer/sso/internal/ports/grpc/resources"
)

type APIServerOption func(srv *APIServer)

func WithAuthManager(authMan *auth.Manager) APIServerOption {
	return func(srv *APIServer) {
		srv.authManager = authMan
	}
}

func WithVersion(version string) APIServerOption {
	return func(srv *APIServer) {
		srv.version = version
	}
}

func WithVerifyManager(verifyMan *resources.Verify) APIServerOption {
	return func(srv *APIServer) {
		srv.verifyMan = verifyMan
	}
}

func WithDatastore(db drivers.DataStore) APIServerOption {
	return func(srv *APIServer) {
		srv.db = db
	}
}
