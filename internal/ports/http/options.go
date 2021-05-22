package http

import (
	"github.com/JetBrainer/sso/internal/domain/manager/auth"
	"github.com/JetBrainer/sso/internal/domain/manager/monitoring"
	"github.com/JetBrainer/sso/pkg/validation"
)

type APIServerOption func(srv *APIServer)

func WithAuthManager(authMan *auth.Manager) APIServerOption {
	return func(srv *APIServer) {
		srv.authManager = authMan
	}
}

func WithMonitoringManager(monitMan *monitoring.Manager) APIServerOption {
	return func(srv *APIServer) {
		srv.monitManager = monitMan
	}
}


func WithValidator(v *validation.Validator) APIServerOption {
	return func(srv *APIServer) {
		srv.validator = v
	}
}

func WithVersion(version string) APIServerOption {
	return func(srv *APIServer) {
		srv.version = version
	}
}
