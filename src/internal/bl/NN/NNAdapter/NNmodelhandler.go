package nn_adapter

import (
	models_dto "annotater/internal/models/dto"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type IModelHandler interface {
	GetModelResp(req ModelRequest) ([]models_dto.Markup, error)
}

type HttpModelHandler struct {
	Url string
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
