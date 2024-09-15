package nn_adapter

import (
	nn_model_handler "annotater/internal/bl/NN/NNAdapter/NNmodelhandler"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	"errors"
	"fmt"
)

var (
	ErrInModelPrediction = errors.New("error in getting predictions")
)

type DetectionModel struct { // NN stands for neural network
	modelHandler nn_model_handler.IModelHandler
}

func NewDetectionModel(handler nn_model_handler.IModelHandler) *DetectionModel {
	return &DetectionModel{modelHandler: handler}
}

func (m *DetectionModel) Predict(document models.DocumentData) ([]models.Markup, error) {
	req := nn_model_handler.ModelRequest{DocumentData: document.DocumentBytes}
	markupsDto, err := m.modelHandler.GetModelResp(req)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("error in getting predictions for document %v", document.ID), err)
	}
	markups := models_dto.FromDtoMarkupSlice(markupsDto)

	return markups, nil

}
