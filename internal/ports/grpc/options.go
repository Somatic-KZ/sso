package grpc

import "github.com/JetBrainer/sso/internal/domain/manager/auth"

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
