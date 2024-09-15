package repo_adapter

import (
	repository "annotater/internal/bl/annotationService/annotattionRepo"
	"annotater/internal/models"
	models_da "annotater/internal/models/modelsDA"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	ErrNothingDelete = errors.New("nothing was deleted")
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
	if tx.Error == gorm.ErrForeignKeyViolated {
		return models.ErrViolatingKeyAnnot
	}

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in adding anotattion")
	}
	return nil
}

func (repo *AnotattionRepositoryAdapter) DeleteAnotattion(id uint64) error { // do we need transactions here?
	tx := repo.db.Where("id = ?", id) //using that because if id is equal to 0 then the first found row will be deleted
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in deleting anotattion")
	}
	fmt.Print(tx.Error)
	tx = tx.Delete(&models_da.Markup{})
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in deleting anotattion")
	}
	/*if tx.RowsAffected == 0 {
		return ErrNothingDelete TODO:: think wether it must be and error
	}*/
	return nil
}

func (repo *AnotattionRepositoryAdapter) GetAnottationByID(id uint64) (*models.Markup, error) {
	var markUpDA models_da.Markup
	tx := repo.db.Where("id = ?", id).First(&markUpDA)

	if tx.Error == gorm.ErrRecordNotFound {
		return nil, models.ErrNotFound
	}

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error in getting anotattion type")
	}
	markUpType, err := models_da.FromDaMarkup(&markUpDA)
	if err != nil {
		return nil, errors.Wrap(err, "Error in getting anotattion type")
	}
	return &markUpType, nil
}
func (repo *AnotattionRepositoryAdapter) GetAnottationsByUserID(id uint64) ([]models.Markup, error) {
	var markUpsDA []models_da.Markup
	tx := repo.db.Where("creator_id = ?", id).Find(&markUpsDA)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error in getting anotattion type")
	}
	markUps, err := models_da.FromDaMarkupSlice(markUpsDA)
	if err != nil {
		return nil, errors.Wrap(err, "Error in getting markups by userID")
	}
	return markUps, err
}

func (repo *AnotattionRepositoryAdapter) GetAllAnottations() ([]models.Markup, error) {
	var markUpsDA []models_da.Markup
	tx := repo.db.Find(&markUpsDA)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error in getting anotattion type")
	}
	markUps, err := models_da.FromDaMarkupSlice(markUpsDA)
	if err != nil {
		return nil, errors.Wrap(err, "Error in getting all markups")
	}
	return markUps, err
}
