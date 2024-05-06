package user_handler

import (
	service "annotater/internal/bl/userService"
	response "annotater/internal/lib/api"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrChangingRole    = errors.New("error changing role")
	ErrDecodingJson    = errors.New("broken request")
	ErrGettingAllUsers = errors.New("error getting all users")
)

type RequestChangeRole struct {
	Login   string      `json:"login"`
	ReqRole models.Role `json:"req_role"`
}

type ResponseGetAllUsers struct {
	response.Response
	Users []models_dto.User `json:"users"`
}

type UserHandler struct {
	logger      *slog.Logger
	userService service.IUserService
}

func NewDocumentHandler(logSrc *slog.Logger, serv service.IUserService) UserHandler {
	return UserHandler{
		logger:      logSrc,
		userService: serv,
	}
}

func (h *UserHandler) ChangeUserPerms() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestChangeRole
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrDecodingJson.Error()))
			h.logger.Error(err.Error())
			return
		}
		err = h.userService.ChangeUserRoleByLogin(req.Login, req.ReqRole)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.logger.Error(err.Error())
			return
		}
		render.JSON(w, r, response.OK())
	}
}

func (h *UserHandler) GetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		users, err := h.userService.GetAllUsers()
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.logger.Error(err.Error())
			return
		}
		usersDTO := models_dto.ToDtoUserSlice(users)
		resp := ResponseGetAllUsers{response.OK(), usersDTO}
		render.JSON(w, r, resp)
	}
}
