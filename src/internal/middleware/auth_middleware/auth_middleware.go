package auth_middleware

import (
	auth_service "annotater/internal/bl/auth"
	response "annotater/internal/lib/api"
	"annotater/internal/models"
	auth_utils "annotater/internal/pkg/authUtils"
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

type contextKeyRole struct{}
type contextKeyID struct{}

func FromIncomingContextRole(ctx context.Context) (models.Role, bool) {
	role, ok := ctx.Value(contextKeyRole{}).(models.Role)
	return role, ok
}

func FromIncomingContextID(ctx context.Context) (uint64, bool) {
	id, ok := ctx.Value(contextKeyID{}).(uint64)
	return id, ok
}

var (
	UserIDContextKey = "contextKeyRole{}"
	RoleContextKey   = "contextKeyID{}"
)

func JwtAuthMiddleware(next http.Handler, secret string, tokenHandler auth_utils.ITokenHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		if token == "" {
			render.JSON(w, r, response.Error("Error in parsing token"))
			render.Status(r, http.StatusBadRequest)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		payload, err := tokenHandler.ParseToken(token, auth_service.SECRET)
		if err != nil {
			if err == auth_utils.ErrParsingToken {
				render.JSON(w, r, response.Error(err.Error()))
				render.Status(r, http.StatusBadRequest)
			} else {
				render.JSON(w, r, response.Error(err.Error()))
				render.Status(r, http.StatusUnauthorized)
			}
			return
		}
		ctx := context.WithValue(r.Context(), UserIDContextKey, payload.ID)
		ctx = context.WithValue(ctx, RoleContextKey, payload.Role)
		//ctx = r.Clone(ctx)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
