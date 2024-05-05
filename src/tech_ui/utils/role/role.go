package role_req

import (
	user_handler "annotater/internal/http-server/handlers/user"
	response "annotater/internal/lib/api"
	"annotater/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var (
	roleChangeUrlPath = "http://localhost:8080/user/role"
)

func ChangeUserRole(client *http.Client, login string, wantedRole models.Role, jwtToken string) error {
	url := roleChangeUrlPath

	jsonReq := user_handler.RequestChangeRole{
		Login:   login,
		ReqRole: wantedRole,
	}
	jsonReqMarshalled, err := json.Marshal(jsonReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonReqMarshalled))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Content-Type", "application/json")

	var respJson *http.Response
	respJson, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	var resp response.Response
	err = render.DecodeJSON(respJson.Body, &resp)
	if err != nil {
		return err
	}
	if resp.Status == response.StatusError {
		return errors.New(resp.Error)
	}
	return nil
}
