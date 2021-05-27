package v1

import (
	"github.com/JetBrainer/sso/internal/domain/manager/auth"
	"github.com/JetBrainer/sso/internal/domain/models"
	"github.com/JetBrainer/sso/pkg/validation"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

type AuthResource struct {
	metrics     *models.Metrics
	authManager *auth.Manager
	validate    *validation.Validator
}

func NewAuth(authMan *auth.Manager, metrics *models.Metrics, validate *validation.Validator) *AuthResource {
	return &AuthResource{
		authManager: authMan,
		validate:    validate,
		metrics:     metrics,
	}
}

func (a AuthResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(a.authManager.TokenAuth()))
		r.Use(NewUserAccessCtx(a.authManager.JWTKey()).ChiMiddleware)

		r.Delete("/signout", a.SignOut)

		r.Put("/verify/reset", a.ResetVerify())
		r.Put("/recovery/reset", a.ResetRecovery())
	})

	r.Group(func(r chi.Router) {
		r.Post("/signin/email", a.SignInByEmail)
		r.Put("/signup", a.SignUP)
		r.Post("/refresh", a.Refresh)
	})

	return r
}
