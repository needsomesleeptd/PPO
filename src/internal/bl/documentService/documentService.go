package service

import (
	nn "annotater/internal/bl/NN"
	repository "annotater/internal/bl/documentService/documentRepo"
	"annotater/internal/models"

	"github.com/pkg/errors"
)

type IDocumentService interface {
	LoadDocument(document models.Document) error
	CheckDocument(document models.Document) ([]*models.Markup, error)
}

type DocumentService struct {
	repo          repository.IDocumentRepository
	neuralNetwork nn.INeuralNetwork
}

func NewDocumentService(pRep repository.IDocumentRepository, pNN nn.INeuralNetwork) IDocumentService {
	return &DocumentService{
		repo:          pRep,
		neuralNetwork: pNN,
	}
}

func (serv *DocumentService) LoadDocument(document models.Document) error {
	/*err := pdf.Validate(bytes.NewReader(document.DocumentData))
	if err != nil {
		errors.Wrap(err, "Loaded file is not a valid pdf file")
	}*/
	err := serv.repo.AddDocument(&document)
	if err != nil {
		return errors.Wrap(err, "Error in loading document")
	}
	return err
}

func (serv *DocumentService) CheckDocument(document models.Document) ([]*models.Markup, error) {
	markups, err := serv.neuralNetwork.Predict(document)
	if err != nil {
		return markups, errors.Wrap(err, "Error in Checking document")
	}
	return markups, err
}
