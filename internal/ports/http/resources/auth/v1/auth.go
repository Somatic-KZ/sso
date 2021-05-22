package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/JetBrainer/sso/internal/domain/manager/auth"
	"github.com/JetBrainer/sso/internal/ports/http/resources"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

var ErrInvalidToken = errors.New("token is incorrect or expired")

type UserAccessCtx struct {
	jwtKey []byte
}

func NewUserAccessCtx(jwtKey []byte) *UserAccessCtx {
	return &UserAccessCtx{
		jwtKey: jwtKey,
	}
}

func (ua UserAccessCtx) ChiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			_ = render.Render(w, r, resources.Unauthorized(err))
			return
		}

		if token == nil || !token.Valid {
			_ = render.Render(w, r, resources.Unauthorized(ErrInvalidToken))
			return
		}

		// инициализируем новый инстанс Claims
		claims := new(auth.Claims)
		tkn, err := jwt.ParseWithClaims(token.Raw, claims, func(token *jwt.Token) (interface{}, error) {
			return ua.jwtKey, nil
		})
		if err != nil {
			_ = render.Render(w, r, resources.BadRequest(ErrInvalidToken))
			return
		}
		if !tkn.Valid {
			_ = render.Render(w, r, resources.Unauthorized(ErrInvalidToken))
			return
		}

		if claims.IsRefresh {
			_ = render.Render(w, r, resources.Unauthorized(ErrInvalidToken))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "tdid", claims.TDID)
		ctx = context.WithValue(ctx, "roles", claims.Roles)

		// токен валидный, пропускаем его
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
