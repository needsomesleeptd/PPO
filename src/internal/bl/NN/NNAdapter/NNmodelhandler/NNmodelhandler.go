package nn_model_handler

import (
	models_dto "annotater/internal/models/dto"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrMarshallingRequest = errors.New("error in Marshalling NN request")
	ErrGettingResponse    = errors.New("error in getting NN response")
)

type IModelHandler interface {
	GetModelResp(req ModelRequest) ([]models_dto.Markup, error)
}

type HttpModelHandler struct {
	Url string
}

func NewHttpModelHandler(url string) IModelHandler {
	return &HttpModelHandler{Url: url}
}

type ModelRequest struct {
	DocumentData []byte `json:"document_data"`
}

func (h *HttpModelHandler) GetModelResp(req ModelRequest) ([]models_dto.Markup, error) {
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Join(ErrMarshallingRequest, err)
	}
	buffer := bytes.NewBuffer(jsonReq)
	jsonResp, err := http.Post(h.Url, "application/json", buffer)
	if err != nil {
		return nil, errors.Join(ErrGettingResponse, err)
	}
	var markupsDto []models_dto.Markup
	json.NewDecoder(jsonResp.Body).Decode(&markupsDto)
	return markupsDto, nil
}
