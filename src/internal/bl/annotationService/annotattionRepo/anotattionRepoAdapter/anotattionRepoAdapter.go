package repo_adapter

import (
	repository "annotater/internal/bl/annotationService/annotattionRepo"
	"annotater/internal/models"
	models_da "annotater/internal/models/modelsDA"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type AnotattionRepositoryAdapter struct {
	db *gorm.DB
}

func NewAnotattionRepositoryAdapter(srcDB *gorm.DB) repository.IAnotattionRepository {
	return &AnotattionRepositoryAdapter{
		db: srcDB,
	}
}

func (repo *AnotattionRepositoryAdapter) AddAnottation(markUp *models.Markup) error {
	markUpDa, err := models_da.ToDaMarkup(*markUp)
	if err != nil {
		return errors.Wrap(err, "Error in getting anotattion type")
	}
	tx := repo.db.Create(markUpDa)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in adding anotattion")
	}
	return nil
}

func (repo *AnotattionRepositoryAdapter) DeleteAnotattion(id uint64) error {
	var markUpDA models_da.Markup
	markUpDA.ID = id
	tx := repo.db.Where("id = ?", id).Delete(markUpDA) //using that because if id is equal to 0 then the first found row will be deleted
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in deleting anotattion")
	}
	return nil
}

func (repo *AnotattionRepositoryAdapter) GetAnottationByID(id uint64) (*models.Markup, error) {
	var markUpDA models_da.Markup
	tx := repo.db.Where("id = ?", id).First(&markUpDA)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error in getting anotattion type")
	}
	markUpType, err := models_da.FromDaMarkup(&markUpDA)
	if err != nil {
		return nil, errors.Wrap(err, "Error in getting anotattion type")
	}
	return &markUpType, nil
}
