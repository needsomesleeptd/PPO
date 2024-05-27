package access_middleware

import (
	service "annotater/internal/bl/userService"
	response "annotater/internal/lib/api"
	"annotater/internal/middleware/auth_middleware"
	"annotater/internal/models"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

var (
	ErrGettingValueFromContext = errors.New("error getting value from context")
	ErrInternalServer          = errors.New("intenal server error")
	ErrAccessDeniedServer      = errors.New("access denied")
)

type AccessMiddleware struct {
	userService service.IUserService
	logger      *logrus.Logger
}

func NewAccessMiddleware(logSrc *logrus.Logger, userServiceSrc service.IUserService) *AccessMiddleware {
	return &AccessMiddleware{userService: userServiceSrc, logger: logSrc}
}

func (ac *AccessMiddleware) AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(auth_middleware.UserIDContextKey).(uint64)
		if !ok {
			ac.logger.WithFields(
				logrus.Fields{
					"src": "AccessMiddleware.AdminOnlyMiddleware"}).
				Info("no valid userID in context")
			render.JSON(w, r, response.Error(ErrInternalServer.Error()))
			render.Status(r, http.StatusBadRequest)
			return
		}

		role, ok := ctx.Value(auth_middleware.RoleContextKey).(models.Role)
		if !ok {
			ac.logger.WithFields(
				logrus.Fields{
					"src":    "AccessMiddleware.AdminOnlyMiddleware",
					"userID": userID}).
				Info("no valid role in context")
			render.JSON(w, r, response.Error(ErrInternalServer.Error()))
			render.Status(r, http.StatusBadRequest)
			return
		}
		if !ac.userService.IsRolePermitted(role, models.Admin) {
			ac.logger.WithFields(
				logrus.Fields{
					"src":    "AccessMiddleware.AdminOnlyMiddleware",
					"userID": userID,
					"role":   role}).
				Info("role is not enough")
			render.JSON(w, r, response.Error(ErrAccessDeniedServer.Error()))
			render.Status(r, http.StatusForbidden)
			return
		}

		ac.logger.WithFields(
			logrus.Fields{
				"src":    "AccessMiddleware.AdminOnlyMiddleware",
				"userID": userID,
				"role":   role}).
			Info("passed admin check")

		next.ServeHTTP(w, r)
	})
}

func (ac *AccessMiddleware) ControllersAndHigherMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(auth_middleware.UserIDContextKey).(uint64)
		if !ok {
			ac.logger.WithFields(
				logrus.Fields{
					"src": "AccessMiddleware.ControllersAndHigherMiddleware"}).
				Info("no valid userID in context")
			render.JSON(w, r, response.Error(ErrInternalServer.Error()))
			render.Status(r, http.StatusBadRequest)
			return
		}

		role, ok := ctx.Value(auth_middleware.RoleContextKey).(models.Role)
		if !ok {
			ac.logger.WithFields(
				logrus.Fields{
					"src":    "AccessMiddleware.ControllersAndHigherMiddleware",
					"userID": userID}).
				Info("no valid role in context")
			render.JSON(w, r, response.Error(ErrInternalServer.Error()))
			render.Status(r, http.StatusBadRequest)
			return
		}
		if !ac.userService.IsRolePermitted(role, models.Controller) {
			ac.logger.WithFields(
				logrus.Fields{
					"src":    "AccessMiddleware.ControllersAndHigherMiddleware",
					"userID": userID,
					"role":   role}).
				Info("role is not enough")
			render.JSON(w, r, response.Error(ErrAccessDeniedServer.Error()))
			render.Status(r, http.StatusForbidden)
			return
		}

		ac.logger.WithFields(
			logrus.Fields{
				"src":    "AccessMiddleware.AdminOnlyMiddleware",
				"userID": userID,
				"role":   role}).
			Info("passed controller and higher check")
		next.ServeHTTP(w, r)
	})
}
