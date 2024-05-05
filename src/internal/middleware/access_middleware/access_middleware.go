package access_middleware

import (
	service "annotater/internal/bl/userService"
	response "annotater/internal/lib/api"
	"annotater/internal/middleware/auth_middleware"
	"annotater/internal/models"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrGettingValueFromContext = errors.New("error getting value from context")
	ErrInternalServer          = errors.New("intenal server error")
	ErrAccessDeniedServer      = errors.New("access denied")
)

type AccessMiddleware struct {
	userService service.IUserService
}

func NewAccessMiddleware(userServiceSrc service.IUserService) *AccessMiddleware {
	return &AccessMiddleware{userService: userServiceSrc}
}

func (ac *AccessMiddleware) AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		role, ok := ctx.Value(auth_middleware.RoleContextKey).(models.Role)
		if !ok {
			render.JSON(w, r, response.Error(ErrInternalServer.Error()))
			render.Status(r, http.StatusBadRequest)
			return
		}
		if !ac.userService.IsRolePermitted(role, models.Admin) {
			render.JSON(w, r, response.Error(ErrAccessDeniedServer.Error()))
			render.Status(r, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (ac *AccessMiddleware) ControllersAndHigherMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		role, ok := ctx.Value(auth_middleware.RoleContextKey).(models.Role)
		if !ok {
			render.JSON(w, r, response.Error(ErrInternalServer.Error()))
			render.Status(r, http.StatusBadRequest)
			return
		}
		if !ac.userService.IsRolePermitted(role, models.Controller) {
			render.JSON(w, r, response.Error(ErrAccessDeniedServer.Error()))
			render.Status(r, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
