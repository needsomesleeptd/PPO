package user_req

import (
	user_handler "annotater/internal/http-server/handlers/user"
	response "annotater/internal/lib/api"
	models_dto "annotater/internal/models/dto"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var (
	gettinAllUsersUrlPath = "http://localhost:8080/user/getUsers"
)

func GetAllUsers(client *http.Client, jwtToken string) ([]models_dto.User, error) {
	url := gettinAllUsersUrlPath

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	var respJson *http.Response
	respJson, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var resp user_handler.ResponseGetAllUsers
	err = render.DecodeJSON(respJson.Body, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Status != response.StatusOK {
		return nil, errors.New(resp.Error)
	}
	return resp.Users, nil
}
