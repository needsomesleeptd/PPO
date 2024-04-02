package nn_adapter

import (
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	"errors"
)

var (
	ErrMarshallingRequest = errors.New("error in Marshalling NN request")
	ErrGettingResponse    = errors.New("error in getting NN response")
	ErrInModelPrediction  = errors.New("error in getting predictions")
)

type DetectionModel struct { // NN stands for neural network
	modelHandler IModelHandler
}

func NewDetectionModel(handler IModelHandler) (*DetectionModel, error) {
	return &DetectionModel{modelHandler: handler}, nil
}

func (m *DetectionModel) Predict(document models.Document) ([]models.Markup, error) {
	req := ModelRequest{DocumentData: document.DocumentData}
	markupsDto, err := m.modelHandler.GetModelResp(req)
	if err != nil {
		return nil, errors.Join(ErrInModelPrediction, err)
	}
	markups := models_dto.FromDtoMarkupSlice(markupsDto)

	return markups, nil

}
