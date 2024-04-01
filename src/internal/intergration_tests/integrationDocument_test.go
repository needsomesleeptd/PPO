package integration_tests

//TODO:: split tests by files

import (
	document_repo_adapter "annotater/internal/bl/documentService/documentRepo/documentRepoAdapter"
	integration_utils "annotater/internal/intergration_tests/utils"
	"annotater/internal/models"
	models_da "annotater/internal/models/modelsDA"
	"bytes"
	"log"
	"testing"

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
	// Open a new database connection for each test
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

//testing Document Service
func (suite *UsecaseRepositoryTestSuite) TestUsecaseAddDocument() {
	var document *models.Document
	userRepo := document_repo_adapter.NewDocumentRepositoryAdapter(suite.db)

	insertedDocument := models.Document{DocumentData: createPDFBuffer(TEST_VALID_PDF)}
	err := userRepo.AddDocument(&insertedDocument)
	suite.Require().NoError(err)
	document, err = userRepo.GetDocumentByID(1)
	suite.Require().NoError(err)
	suite.Assert().Equal(document.DocumentData, insertedDocument.DocumentData)
	suite.Assert().Equal(document.ID, uint64(1))
	// interestingly time format has changed from local to UTC
}

func (suite *UsecaseRepositoryTestSuite) TestUsecaseAddDocumentID() {
	var document *models.Document
	userRepo := document_repo_adapter.NewDocumentRepositoryAdapter(suite.db)
	id := uint64(2)
	insertedDocument := models.Document{ID: id, DocumentData: createPDFBuffer(TEST_VALID_PDF)}
	err := userRepo.AddDocument(&insertedDocument)
	suite.Require().NoError(err)
	document, err = userRepo.GetDocumentByID(id)
	suite.Require().NoError(err)
	suite.Assert().Equal(document.DocumentData, insertedDocument.DocumentData)
	suite.Assert().Equal(document.ID, id)
}

func (suite *UsecaseRepositoryTestSuite) TestUsecaseDeleteDocumentID() {
	document := models.Document{}
	userRepo := document_repo_adapter.NewDocumentRepositoryAdapter(suite.db)
	id := uint64(2)
	insertedDocument := models.Document{ID: id, DocumentData: createPDFBuffer(TEST_VALID_PDF)}
	err := userRepo.AddDocument(&insertedDocument)
	suite.Require().NoError(err)
	suite.Assert().NoError(suite.db.Table("documents").First(&document, models.Document{ID: id}).Error)
	err = userRepo.DeleteDocumentByID(id)
	suite.Require().NoError(err)
	suite.Assert().Error(suite.db.Table("documents").First(&document, models.Document{ID: id}).Error)

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(UsecaseRepositoryTestSuite))
}
