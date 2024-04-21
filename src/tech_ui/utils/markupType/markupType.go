package annot_type_req

import (
	annot_type_handler "annotater/internal/http-server/handlers/annotType"
	response "annotater/internal/lib/api"
	models_dto "annotater/internal/models/dto"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var (
	annotTypesUrlPath = "http://localhost:8080/annotType/"
)

func GetMarkupTypesCreatorID(client *http.Client, jwtToken string) ([]models_dto.MarkupType, error) {
	url := annotTypesUrlPath + "creatorID"

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

	var resp annot_type_handler.ResponseGetByIDs
	err = render.DecodeJSON(respJson.Body, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Response != response.OK() {
		return nil, errors.New(resp.Error)
	}
	return resp.MarkupTypes, nil
}

func GetMarkupTypeByCreatorID(client *http.Client, labelName string, description string, jwtToken string) error {
	url := annotTypesUrlPath + "add"

	reqBody, err := json.Marshal(annot_type_handler.RequestAnnotType{Description: description, ClassName: labelName})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)
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
	if resp != response.OK() {
		return errors.New(resp.Error)
	}
	return nil
}
