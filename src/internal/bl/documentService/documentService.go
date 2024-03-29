package service

import (
	nn "annotater/internal/bl/NN"
	repository "annotater/internal/bl/documentService/documentRepo"
	"annotater/internal/models"
	"bytes"

	"github.com/pkg/errors"
	"github.com/telkomdev/go-filesig"
)

const (
	LOADING_DOCUMENT_ERR_STR  = "Error in loading document"
	CHECKING_DOCUMENT_ERR_STR = "Error in Checking document"
	DOCUMENT_FORMAT_ERR_STR   = "Error document loaded in wrong format"
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
	isValid := filesig.IsPdf(bytes.NewReader(document.DocumentData))
	if !isValid {
		return errors.New(DOCUMENT_FORMAT_ERR_STR)
	}
	err := serv.repo.AddDocument(&document)
	if err != nil {
		return errors.Wrap(err, LOADING_DOCUMENT_ERR_STR)
	}
	return err
}

func (serv *DocumentService) CheckDocument(document models.Document) ([]*models.Markup, error) {

	isValid := filesig.IsPdf(bytes.NewReader(document.DocumentData))
	if !isValid {
		return nil, errors.New(DOCUMENT_FORMAT_ERR_STR)
	}
	markups, err := serv.neuralNetwork.Predict(document)
	if err != nil {
		return markups, errors.Wrap(err, CHECKING_DOCUMENT_ERR_STR)
	}
	return markups, err
}
