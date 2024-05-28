package user_handler

import (
	service "annotater/internal/bl/userService"
	response "annotater/internal/lib/api"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
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
	logger      *logrus.Logger
	userService service.IUserService
}

func NewDocumentHandler(logSrc *logrus.Logger, serv service.IUserService) UserHandler {
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
			h.logger.Warn(err.Error())
			return
		}
		err = h.userService.ChangeUserRoleByLogin(req.Login, req.ReqRole)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.logger.Warn(err.Error())
			return
		}
		h.logger.Infof("successfully changed role of user with login %v  to role %v\n", req.Login, req.ReqRole)
		render.JSON(w, r, response.OK())
	}
}

func (h *UserHandler) GetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		users, err := h.userService.GetAllUsers()
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.logger.Warn(err.Error())
			return
		}
		usersDTO := models_dto.ToDtoUserSlice(users)
		resp := ResponseGetAllUsers{response.OK(), usersDTO}
		h.logger.Infof("succesfully got all users\n")
		render.JSON(w, r, resp)
	}
}
