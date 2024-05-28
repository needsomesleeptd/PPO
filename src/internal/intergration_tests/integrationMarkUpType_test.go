package integration_tests

import (
	service "annotater/internal/bl/anotattionTypeService"
	repo_adapter "annotater/internal/bl/anotattionTypeService/anottationTypeRepo/anotattionTypeRepoAdapter"
	integration_utils "annotater/internal/intergration_tests/utils"
	"annotater/internal/models"
	models_da "annotater/internal/models/modelsDA"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MarkupTypeTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *MarkupTypeTestSuite) SetupTest() {
	db, err := gorm.Open(postgres.New(integration_utils.TestCfg), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models_da.MarkupType{})
	suite.db = db
}

func (suite *MarkupTypeTestSuite) TearDownTest() {
	suite.db.Migrator().DropTable(&models_da.MarkupType{})
}
func (suite *MarkupTypeTestSuite) TestUsecaseAddMarkUpType() {
	anotattionTypeRepo := repo_adapter.NewAnotattionTypeRepositoryAdapter(suite.db)
	anotattionTypeService := service.NewAnotattionTypeService(anotattionTypeRepo)
	id := uint64(1)
	markUpType := models.MarkupType{ID: id}
	gotMarkUpType := models_da.MarkupType{ID: id}
	suite.Require().Error(suite.db.Model(&models_da.MarkupType{}).Where("id = ?", id).Take(&gotMarkUpType).Error)
	err := anotattionTypeService.AddAnottationType(&markUpType)
	suite.Require().NoError(err)
	suite.Assert().NoError(suite.db.Model(&models_da.MarkupType{}).Where("id = ?", id).Take(&gotMarkUpType).Error)
	suite.Assert().Equal(models_da.FromDaMarkupType(&gotMarkUpType), markUpType)
}

func (suite *MarkupTypeTestSuite) TestUsecaseGetMarkUpType() {
	anotattionTypeRepo := repo_adapter.NewAnotattionTypeRepositoryAdapter(suite.db)
	anotattionTypeService := service.NewAnotattionTypeService(anotattionTypeRepo)
	id := uint64(1)
	markUpTypeDa := models_da.MarkupType{ID: id, CreatorID: 12}
	suite.Require().NoError(suite.db.Create(&markUpTypeDa).Error)

	markUpType, err := anotattionTypeService.GetAnottationTypeByID(id)
	suite.Require().NoError(err)
	suite.Require().Equal(*markUpType, models_da.FromDaMarkupType(&markUpTypeDa))
}

func (suite *MarkupTypeTestSuite) TestUsecaseDeleteMarkUpType() {
	anotattionRepo := repo_adapter.NewAnotattionTypeRepositoryAdapter(suite.db)
	anotattionService := service.NewAnotattionTypeService(anotattionRepo)
	id := uint64(1)
	markUpType := models.MarkupType{ID: id, CreatorID: 12}
	gotMarkUp := models_da.MarkupType{ID: id}
	err := anotattionService.AddAnottationType(&markUpType)
	suite.Require().NoError(err)
	suite.Require().NoError(suite.db.Model(&models_da.MarkupType{}).Where("id = ?", id).Take(&gotMarkUp).Error)
	err = anotattionService.DeleteAnotattionType(id)
	suite.Require().NoError(err)
	suite.Require().Error(suite.db.Model(&models_da.MarkupType{}).Where("id = ?", id).Take(&gotMarkUp).Error)
}

func TestSuiteMarkupType(t *testing.T) {
	suite.Run(t, new(MarkupTypeTestSuite))
}
