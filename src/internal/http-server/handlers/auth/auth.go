package auth_handler

import (
	auth_service "annotater/internal/bl/auth"
	response "annotater/internal/lib/api"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrDecodingJson = errors.New("broken request")
	ErrSignUp       = errors.New("signup error")
	ErrSignIn       = errors.New("signin error")
)

const (
	COOKIE_NAME = "auth_jwt"
)

type RequestSignUp struct {
	User models_dto.User `json:"user"`
}
type RequestSignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ResponseSignIn struct {
	Response response.Response
	Jwt      string `json:"jwt"`
	ID       uint64 `json:"userID"`
}

func SignUp(authService auth_service.IAuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestSignUp
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrDecodingJson.Error())) //TODO:: add logging here
			return
		}

		candidate := models_dto.FromDtoUser(&req.User)
		err = authService.SignUp(&candidate)
		if err != nil {
			render.JSON(w, r, response.Error(ErrSignUp.Error()))
			return
		}

		render.JSON(w, r, response.OK())
	}
}

func SignIn(authService auth_service.IAuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestSignIn
		var tokenStr string
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrDecodingJson.Error()))
			return
		}
		candidate := models.User{Login: req.Login, Password: req.Password}
		tokenStr, err = authService.SignIn(&candidate)

		if err != nil {
			render.JSON(w, r, response.Error(ErrSignIn.Error()))
			return
		}

		resp := ResponseSignIn{Response: response.OK(), Jwt: tokenStr, ID: candidate.ID}

		render.JSON(w, r, resp)
	}
}