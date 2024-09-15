package auth_ui

import (
	auth_handler "annotater/internal/http-server/handlers/auth"
	response "annotater/internal/lib/api"
	models_dto "annotater/internal/models/dto"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var (
	authPath = "http://localhost:8080/user/"
)

func AddJwtHeader(req *http.Request, jwtToken string) {

	req.Header.Add("Authorization", "Bearer "+jwtToken)
}

func SignIn(client *http.Client, login string, password string) (string, error) {

	url := authPath + "SignIn"

	reqBody := auth_handler.RequestSignIn{Login: login, Password: password}
	reqBodyJson, _ := json.Marshal(reqBody)
	respGot, err := http.Post(url, "application/json", bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return "", err
	}
	var resp auth_handler.ResponseSignIn
	err = render.DecodeJSON(respGot.Body, &resp)

	if err != nil {
		return "", err
	}
	if resp.Response.Status == response.StatusError {
		return "", errors.New(resp.Response.Error)
	}
	return resp.Jwt, nil
}

func SignUp(client *http.Client, login string, password string) (string, error) {

	url := authPath + "SignUp"

	user := models_dto.User{Login: login, Password: password}
	reqBody := auth_handler.RequestSignUp{User: user}
	reqBodyJson, _ := json.Marshal(reqBody)
	respGot, err := http.Post(url, "application/json", bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return "", err
	}
	resp := response.Response{}
	err = render.DecodeJSON(respGot.Body, &resp)
	if err != nil {
		return "", err
	}
	if resp.Status == response.StatusOK {
		return response.StatusOK, nil
	} else {
		return "", errors.New(resp.Error)
	}
}
