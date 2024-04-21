package service

import (
	nn "annotater/internal/bl/NN"
	annot_type_repository "annotater/internal/bl/anotattionTypeService/anottationTypeRepo"
	dock_repository "annotater/internal/bl/documentService/documentRepo"
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
	CheckDocument(document models.Document) ([]models.Markup, []models.MarkupType, error)
}

type DocumentService struct {
	docRepo       dock_repository.IDocumentRepository
	annotTypeRepo annot_type_repository.IAnotattionTypeRepository
	neuralNetwork nn.INeuralNetwork
}

func NewDocumentService(docRep dock_repository.IDocumentRepository, pNN nn.INeuralNetwork, typeRep annot_type_repository.IAnotattionTypeRepository) IDocumentService {
	return &DocumentService{
		docRepo:       docRep,
		neuralNetwork: pNN,
		annotTypeRepo: typeRep,
	}
}

func (serv *DocumentService) LoadDocument(document models.Document) error {
	isValid := filesig.IsPdf(bytes.NewReader(document.DocumentData))
	if !isValid {
		return errors.New(DOCUMENT_FORMAT_ERR_STR)
	}
	err := serv.docRepo.AddDocument(&document)
	if err != nil {
		return errors.Wrap(err, LOADING_DOCUMENT_ERR_STR)
	}
	return err
}

func (serv *DocumentService) CheckDocument(document models.Document) ([]models.Markup, []models.MarkupType, error) {

	isValid := filesig.IsPdf(bytes.NewReader(document.DocumentData))
	if !isValid {
		return nil, nil, errors.New(DOCUMENT_FORMAT_ERR_STR)
	}
	markups, err := serv.neuralNetwork.Predict(document)
	if err != nil {
		return nil, nil, errors.Wrap(err, CHECKING_DOCUMENT_ERR_STR)
	}
	ids := make([]uint64, len(markups))
	for i := range ids {
		ids[i] = markups[i].ClassLabel
	}
	markupTypes, err := serv.annotTypeRepo.GetAnottationTypesByIDs(ids)

	if len(markupTypes) != len(markups) {
		return nil, nil, errors.Wrap(err, CHECKING_DOCUMENT_ERR_STR)
	}
	if err != nil {
		return nil, nil, errors.Wrap(err, CHECKING_DOCUMENT_ERR_STR)
	}
	return markups, markupTypes, err
}
