package user_handler

import (
	service "annotater/internal/bl/userService"
	response "annotater/internal/lib/api"
	"annotater/internal/models"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrChangingRole = errors.New("error changing role")
	ErrDecodingJson = errors.New("broken request")
)

type RequestChangeRole struct {
	Login   string      `json:"login"`
	ReqRole models.Role `json:"req_role"`
}

func ChangeUserPerms(userService service.IUserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestChangeRole
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrDecodingJson.Error()))
			return
		}
		err = userService.ChangeUserRoleByLogin(req.Login, req.ReqRole)
		if err != nil {
			render.JSON(w, r, response.Error(ErrChangingRole.Error()))
			return
		}
		render.JSON(w, r, response.OK())
	}
}
