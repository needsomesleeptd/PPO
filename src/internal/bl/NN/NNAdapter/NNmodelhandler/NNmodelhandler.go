package nn_model_handler

import (
	models_dto "annotater/internal/models/dto"
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
)

var (
	ErrMarshallingRequest = errors.New("error in Marshalling NN request")
	ErrGettingResponse    = errors.New("error in getting NN response")
	ErrCreatingFormData   = errors.New("error in creating Form Data")
	ErrCreatingRequest    = errors.New("error in creating request")

	PdfFieldName = "document_data"
	PdfFileName  = "request_file"
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
	/*jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Join(ErrMarshallingRequest, err)
	}*/

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(PdfFieldName, PdfFileName)
	if err != nil {
		return nil, errors.Join(ErrCreatingFormData, err)
	}

	_, err = part.Write(req.DocumentData)
	if err != nil {
		return nil, errors.Join(ErrCreatingFormData, err)
	}
	reqModel, err := http.NewRequest("POST", "application/json", body)

	if err != nil {
		return nil, errors.Join(ErrCreatingFormData, err)
	}
	reqModel.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{} // TODO:: pass the client from outside
	jsonResp, err := client.Do(reqModel)
	if err != nil {
		return nil, errors.Join(ErrGettingResponse, err)
	}

	var markupsDto []models_dto.Markup
	json.NewDecoder(jsonResp.Body).Decode(&markupsDto)
	return markupsDto, nil
}
