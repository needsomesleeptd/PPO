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
	ErrMarshallingRequest    = errors.New("error in Marshalling NN request")
	ErrUnMarshallingResponse = errors.New("error in Unmarshalling NN request")
	ErrGettingResponse       = errors.New("error in getting NN response")
	ErrCreatingFormData      = errors.New("error in creating Form Data")
	ErrCreatingRequest       = errors.New("error in creating request")

	PdfFieldName = "document_data"
	PdfFileName  = "request_file.pdf"
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
	writer.Close()
	reqModel, err := http.NewRequest("POST", h.Url, body)

	if err != nil {
		return nil, errors.Join(ErrCreatingFormData, err)
	}
	reqModel.Header.Set("Content-Type", writer.FormDataContentType())

	jsonResp, err := http.DefaultClient.Do(reqModel)
	if err != nil {
		return nil, errors.Join(ErrGettingResponse, err)
	}

	var markupsDto []models_dto.Markup

	err = json.NewDecoder(jsonResp.Body).Decode(&markupsDto)
	if err != nil {
		return nil, errors.Join(ErrUnMarshallingResponse, err)
	}
	return markupsDto, nil
}
