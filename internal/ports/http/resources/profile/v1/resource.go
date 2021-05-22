package v1

import (
	"github.com/JetBrainer/sso/internal/domain/manager/auth"
	"github.com/JetBrainer/sso/internal/domain/models"
	"github.com/JetBrainer/sso/internal/domain/models/api"
	"github.com/JetBrainer/sso/internal/ports/http/resources/auth/v1"
	"github.com/JetBrainer/sso/pkg/validation"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

type ProfileResource struct {
	authManager  *auth.Manager
	validate     *validation.Validator
}

func NewProfile(authMan *auth.Manager, validate *validation.Validator) *ProfileResource {
	return &ProfileResource{
		authManager:  authMan,
		validate:     validate,
	}
}

func (p ProfileResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(p.authManager.TokenAuth()))
		r.Use(v1.NewUserAccessCtx(p.authManager.JWTKey()).ChiMiddleware)

		r.Get("/", p.Profile)
		r.Put("/", p.ProfileUpdate)

		r.Put("/password", p.UpdatePassword)

		r.Get("/receivers", p.Receivers)
		r.Get("/addresses", p.ReceiversAddress)
	})

	return r
}

func userProfileRequestMerge(request *api.ProfileUpdateRequest, user *models.User) *models.User {
	if request.FirstName != nil {
		user.FirstName = *request.FirstName
	}

	if request.LastName != nil {
		user.LastName = *request.LastName
	}

	if request.Patronymic != nil {
		user.Patronymic = *request.Patronymic
	}

	if request.Email != nil {
		user.Email = *request.Email
	}

	if request.Language != nil {
		user.Language = *request.Language
	}

	if request.Sex != nil {
		user.Sex = *request.Sex
	}

	if request.IIN != nil {
		user.IIN = *request.IIN
	}
	return user
}

func profileKindByUser(user models.User) string {
	if user.Password == "" {
		return "fast"
	}
	return "full"
}
