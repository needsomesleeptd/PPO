package integration_tests

//TODO:: split tests by files

import (
	nn_adapter "annotater/internal/bl/NN/NNAdapter"
	nn_model_handler "annotater/internal/bl/NN/NNAdapter/NNmodelhandler"
	service "annotater/internal/bl/documentService"
	document_repo_adapter "annotater/internal/bl/documentService/documentRepo/documentRepoAdapter"
	integration_utils "annotater/internal/intergration_tests/utils"
	mock_nn_model_handler "annotater/internal/mocks/bl/NN/NNAdapter/NNmodelhandler"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	models_da "annotater/internal/models/modelsDA"
	"bytes"
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/signintech/gopdf"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	TEST_VALID_PDF *gopdf.GoPdf = &gopdf.GoPdf{}
)

func createPDFBuffer(pdf *gopdf.GoPdf) []byte {
	if pdf == nil {
		return []byte{1}
	}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	var buf bytes.Buffer
	pdf.WriteTo(&buf)

	return buf.Bytes()
}

type UsecaseRepositoryTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *UsecaseRepositoryTestSuite) SetupTest() {
	db, err := gorm.Open(postgres.New(integration_utils.TestCfg), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models_da.Document{})

	suite.db = db
}

func (suite *UsecaseRepositoryTestSuite) TearDownTest() {

	suite.db.Migrator().DropTable(&models_da.Document{})

}

// testing Document Service
func (suite *UsecaseRepositoryTestSuite) TestUsecaseAddDocument() {
	var document *models.DocumentMetaData
	userRepo := document_repo_adapter.NewDocumentRepositoryAdapter(suite.db)
	id := uuid.UUID{2}
	insertedDocument := models.DocumentMetaData{ID: id, DocumentData: createPDFBuffer(TEST_VALID_PDF)}
	err := userRepo.AddDocument(&insertedDocument)
	suite.Require().NoError(err)
	document, err = userRepo.GetDocumentByID(id)
	suite.Require().NoError(err)
	suite.Assert().Equal(document.DocumentData, insertedDocument.DocumentData)
	suite.Assert().Equal(document.ID, id)

}

func (suite *UsecaseRepositoryTestSuite) TestUsecaseLoadDocument() {
	var document *models.DocumentMetaData
	userRepo := document_repo_adapter.NewDocumentRepositoryAdapter(suite.db)
	handler := mock_nn_model_handler.NewMockIModelHandler(&gomock.Controller{})
	nn := nn_adapter.NewDetectionModel(handler)
	service := service.NewDocumentService(userRepo, nn)
	id := uuid.UUID{2}
	insertedDocument := models.DocumentMetaData{ID: id, DocumentData: createPDFBuffer(TEST_VALID_PDF)}
	err := service.LoadDocument(insertedDocument)
	suite.Assert().NoError(err)
	document, err = userRepo.GetDocumentByID(id)
	suite.Require().NoError(err)
	suite.Assert().Equal(document.DocumentData, insertedDocument.DocumentData)
	suite.Assert().Equal(document.ID, id)
}

func (suite *UsecaseRepositoryTestSuite) TestUsecaseCheckDocument() {

	userRepo := document_repo_adapter.NewDocumentRepositoryAdapter(suite.db)
	ctrl := gomock.NewController(suite.T())
	handler := mock_nn_model_handler.NewMockIModelHandler(ctrl)
	nn := nn_adapter.NewDetectionModel(handler)
	service := service.NewDocumentService(userRepo, nn)
	id := uuid.UUID{2}
	insertedDocument := models.DocumentMetaData{ID: id, DocumentData: createPDFBuffer(TEST_VALID_PDF)}
	marups := []models_dto.Markup{
		{ErrorBB: []float32{0.1, 0.2, 0.3, 0.2}, ClassLabel: 1},
		{ErrorBB: []float32{0.3, 0.2, 0.1, 0.3}, ClassLabel: 2},
	}
	req := nn_model_handler.ModelRequest{DocumentData: insertedDocument.DocumentData}
	handler.EXPECT().GetModelResp(req).Return(marups, nil)
	res, err := service.CheckDocument(insertedDocument)
	suite.Assert().NoError(err)
	suite.Assert().Equal(res, models_dto.FromDtoMarkupSlice(marups))

}

func (suite *UsecaseRepositoryTestSuite) TestUsecaseCheckDocument() {

	userRepo := document_repo_adapter.NewDocumentRepositoryAdapter(suite.db)
	ctrl := gomock.NewController(suite.T())
	handler := mock_nn_model_handler.NewMockIModelHandler(ctrl)
	nn := nn_adapter.NewDetectionModel(handler)
	service := service.NewDocumentService(userRepo, nn)
	id := uint64(2)
	insertedDocument := models.Document{ID: id, DocumentData: createPDFBuffer(TEST_VALID_PDF)}
	marups := []models_dto.Markup{
		{ErrorBB: []float32{0.1, 0.2, 0.3, 0.2}, ClassLabel: 1},
		{ErrorBB: []float32{0.3, 0.2, 0.1, 0.3}, ClassLabel: 2},
	}
	req := nn_model_handler.ModelRequest{DocumentData: insertedDocument.DocumentData}
	handler.EXPECT().GetModelResp(req).Return(marups, nil)
	res, err := service.CheckDocument(insertedDocument)
	suite.Assert().NoError(err)
	suite.Assert().Equal(res, models_dto.FromDtoMarkupSlice(marups))

}

func (suite *UsecaseRepositoryTestSuite) TestUsecaseDeleteDocumentID() {
	document := models.DocumentMetaData{}
	userRepo := document_repo_adapter.NewDocumentRepositoryAdapter(suite.db)
	id := uuid.UUID{2}
	insertedDocument := models.DocumentMetaData{ID: id, DocumentData: createPDFBuffer(TEST_VALID_PDF)}
	err := userRepo.AddDocument(&insertedDocument)
	suite.Require().NoError(err)
	suite.Assert().NoError(suite.db.Table("documents").First(&document, models.DocumentMetaData{ID: id}).Error)
	err = userRepo.DeleteDocumentByID(id)
	suite.Require().NoError(err)
	suite.Assert().Error(suite.db.Table("documents").First(&document, models.DocumentMetaData{ID: id}).Error)

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(UsecaseRepositoryTestSuite))
}
