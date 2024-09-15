package integration_tests

import (
	annot_service "annotater/internal/bl/annotationService"
	annot_repo_adapter "annotater/internal/bl/annotationService/annotattionRepo/anotattionRepoAdapter"
	integration_utils "annotater/internal/intergration_tests/utils"
	"annotater/internal/models"
	models_da "annotater/internal/models/modelsDA"
	"bytes"
	"image"
	"image/png"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var TEST_VALID_PNG_IMG *image.RGBA = image.NewRGBA(image.Rect(0, 0, 100, 100))

func createPNGBuffer(img *image.RGBA) []byte {
	if img == nil {
		return nil
	}
	pngBuf := new(bytes.Buffer)
	png.Encode(pngBuf, img)
	return pngBuf.Bytes()
}

type MarkupTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *MarkupTestSuite) SetupTest() {
	db, err := gorm.Open(postgres.New(integration_utils.TestCfg), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models_da.Markup{})
	suite.db = db
}

func (suite *MarkupTestSuite) TestUsecaseAddMarkUp() {
	anotattionRepo := annot_repo_adapter.NewAnotattionRepositoryAdapter(suite.db)
	anotattionService := annot_service.NewAnnotattionService(anotattionRepo)
	id := uint64(1)
	markUp := models.Markup{ID: id, PageData: createPNGBuffer(TEST_VALID_PNG_IMG)}
	gotMarkUp := models_da.Markup{ID: id}
	suite.Require().Error(suite.db.Model(&models_da.Markup{}).Where("id = ?", id).Take(&gotMarkUp).Error)
	err := anotattionService.AddAnottation(&markUp)
	suite.Require().NoError(err)
	suite.Assert().NoError(suite.db.Model(&models_da.Markup{}).Where("id = ?", id).Take(&gotMarkUp).Error)
	markUpNew, _ := models_da.FromDaMarkup(&gotMarkUp)
	suite.Assert().Equal(markUpNew, markUp)
}

func (suite *MarkupTestSuite) TestUsecaseDeleteMarkUp() {
	anotattionRepo := annot_repo_adapter.NewAnotattionRepositoryAdapter(suite.db)
	anotattionService := annot_service.NewAnnotattionService(anotattionRepo)
	id := uint64(1)
	markUp := models.Markup{ID: id, PageData: createPNGBuffer(TEST_VALID_PNG_IMG)}
	gotMarkUp := models_da.Markup{ID: id}
	err := anotattionService.AddAnottation(&markUp)
	suite.Require().NoError(err)
	suite.Require().NoError(suite.db.Model(&models_da.Markup{}).Where("id = ?", id).Take(&gotMarkUp).Error)
	err = anotattionService.DeleteAnotattion(id)
	suite.Require().NoError(err)
	suite.Require().Error(suite.db.Model(&models_da.Markup{}).Where("id = ?", id).Take(&gotMarkUp).Error)
}

func (suite *MarkupTestSuite) TestUsecaseGetMarkUp() {
	anotattionRepo := annot_repo_adapter.NewAnotattionRepositoryAdapter(suite.db)
	anotattionService := annot_service.NewAnnotattionService(anotattionRepo)
	id := uint64(1)
	markUpDa := models_da.Markup{ID: id, PageData: createPNGBuffer(TEST_VALID_PNG_IMG)}
	suite.Require().NoError(suite.db.Create(&markUpDa).Error)

	markUp, err := anotattionService.GetAnottationByID(id)
	suite.Require().NoError(err)
	markUpNew, _ := models_da.FromDaMarkup(&markUpDa)
	suite.Require().Equal(*markUp, markUpNew)
}

func (suite *MarkupTestSuite) TearDownTest() {
	suite.db.Migrator().DropTable(&models_da.Markup{})
}

func TestSuiteMarkup(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
