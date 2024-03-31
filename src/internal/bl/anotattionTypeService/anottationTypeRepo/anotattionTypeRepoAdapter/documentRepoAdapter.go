package repo_adapter

import (
	repository "annotater/internal/bl/anotattionTypeService/anottationTypeRepo"
	"annotater/internal/models"
	models_da "annotater/internal/models/modelsDA"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type AnotattionTypeRepositoryAdapter struct {
	db *gorm.DB
}

func NewAnotattionTypeRepositoryAdapter(srcDB *gorm.DB) repository.IAnotattionTypeRepository {
	return &AnotattionTypeRepositoryAdapter{
		db: srcDB,
	}
}

func (repo *AnotattionTypeRepositoryAdapter) AddAnottationType(markUp *models.MarkupType) error {
	tx := repo.db.Create(models_da.ToDaMarkupType(*markUp))
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in adding anotattion type")
	}
	return nil
}

func (repo *AnotattionTypeRepositoryAdapter) DeleteAnotattionType(id uint64) error {
	var markUpTypeDA models_da.MarkupType
	markUpTypeDA.ID = id
	tx := repo.db.Delete(markUpTypeDA)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in deleting anotattion type")
	}
	return nil
}

func (repo *AnotattionTypeRepositoryAdapter) GetAnottationTypeByID(id uint64) (*models.MarkupType, error) {
	var markUpTypeDA models_da.MarkupType
	markUpTypeDA.ID = id
	tx := repo.db.First(markUpTypeDA)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error in getting anotattion type")
	}
	markUpType := models_da.FromDaMarkupType(&markUpTypeDA)
	return &markUpType, nil
}
