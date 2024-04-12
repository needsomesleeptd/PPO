package auth_middleware

import (
	auth_service "annotater/internal/bl/auth"
	response "annotater/internal/lib/api"
	auth_utils "annotater/internal/pkg/authUtils"
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

const UserIDContextKey = "userID"

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

		ctx := context.WithValue(r.Context(), UserIDContextKey, payload.ID) //TODO:: find out why no strings

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
