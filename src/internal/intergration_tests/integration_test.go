package integration_tests

//TODO:: split tests by files

import (
	annot_service "annotater/internal/bl/annotationService"
	annot_repo_adapter "annotater/internal/bl/annotationService/annotattionRepo/anotattionRepoAdapter"
	service "annotater/internal/bl/anotattionTypeService"
	repo_adapter "annotater/internal/bl/anotattionTypeService/anottationTypeRepo/anotattionTypeRepoAdapter"
	auth_service "annotater/internal/bl/auth"
	document_repo_adapter "annotater/internal/bl/documentService/documentRepo/documentRepoAdapter"
	user_repo_adapter "annotater/internal/bl/userService/userRepo/userRepoAdapter"
	"annotater/internal/models"
	models_da "annotater/internal/models/modelsDA"
	auth_utils "annotater/internal/pkg/authUtils"
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"testing"

	"github.com/signintech/gopdf"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	CONN_STR                        = "host=localhost user=andrew password=1 database=lab01db port=5432"
	testCfg                         = postgres.Config{DSN: CONN_STR}
	TEST_VALID_PNG_IMG *image.RGBA  = image.NewRGBA(image.Rect(0, 0, 100, 100))
	TEST_VALID_PDF     *gopdf.GoPdf = &gopdf.GoPdf{}
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

func createPNGBuffer(img *image.RGBA) []byte {
	if img == nil {
		return nil
	}
	pngBuf := new(bytes.Buffer)
	png.Encode(pngBuf, img)
	return pngBuf.Bytes()
}

type UsecaseRepositoryTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *UsecaseRepositoryTestSuite) SetupTest() {
	// Open a new database connection for each test
	db, err := gorm.Open(postgres.New(testCfg), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// Automatically migrate the schema for each test
	db.AutoMigrate(&models_da.User{})
	db.AutoMigrate(&models_da.Document{})
	db.AutoMigrate(&models_da.Markup{})
	db.AutoMigrate(&models_da.MarkupType{})
	suite.db = db
}

func (suite *UsecaseRepositoryTestSuite) TearDownTest() {
	// Delete the test table after each test
	suite.db.Migrator().DropTable(&models_da.User{})
	suite.db.Migrator().DropTable(&models_da.Document{})
	suite.db.Migrator().DropTable(&models_da.Markup{})
	suite.db.Migrator().DropTable(&models_da.MarkupType{})
}

//auth tests
func (suite *UsecaseRepositoryTestSuite) TestUsecaseSignUp() {
	userRepo := user_repo_adapter.NewUserRepositoryAdapter(suite.db)
	hasher := auth_utils.NewPasswordHashCrypto()
	tokenHandler := auth_utils.NewJWTTokenHandler()
	key := "key"
	userService := auth_service.NewAuthService(userRepo, hasher, tokenHandler, key)
	var id uint64 = 1
	user := models.User{
		ID:       id,
		Login:    "test_user",
		Password: "test_password",
		Name:     "Test",
		Surname:  "User",
		Role:     models.Admin,
		Group:    "test_group",
	}
	var gotUser *models.User
	err := userService.SignUp(&user)
	suite.Require().NoError(err)
	gotUser, err = userRepo.GetUserByID(id)
	suite.Require().NoError(err)
	fmt.Print(user, gotUser)
	suite.Require().NoError(hasher.ComparePasswordhash(user.Password, gotUser.Password))
	suite.Require().NoError(err)

	var gotUserDa *models_da.User
	suite.Require().NoError(suite.db.Model(&models_da.User{}).Where("id = ?", id).Take(&gotUserDa).Error)
	suite.Assert().Equal(*gotUser, models_da.FromDaUser(gotUserDa))
}

func (suite *UsecaseRepositoryTestSuite) TestUsecaseSignIn() {
	userRepo := user_repo_adapter.NewUserRepositoryAdapter(suite.db)
	hasher := auth_utils.NewPasswordHashCrypto()
	tokenHandler := auth_utils.NewJWTTokenHandler()
	key := "key"
	userService := auth_service.NewAuthService(userRepo, hasher, tokenHandler, key)
	var id uint64 = 1
	user := models.User{
		ID:       id,
		Login:    "test_user",
		Password: "test_password",
		Name:     "Test",
		Surname:  "User",
		Role:     models.Admin,
		Group:    "test_group",
	}
	err := userService.SignUp(&user)
	suite.Require().NoError(err)
	token, err := userService.SignIn(&user)
	suite.Require().NoError(err)

	err = tokenHandler.ValidateToken(token, key)
	suite.Require().NoError(err)

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

// Testing adding markUp
func (suite *UsecaseRepositoryTestSuite) TestUsecaseAddMarkUp() {
	anotattionRepo := annot_repo_adapter.NewAnotattionRepositoryAdapter(suite.db)
	anotattionService := annot_service.NewAnnotattionService(anotattionRepo)
	id := uint64(1)
	markUp := models.Markup{ID: id, PageData: createPNGBuffer(TEST_VALID_PNG_IMG)}
	gotMarkUp := models_da.Markup{ID: id}
	suite.Require().Error(suite.db.Model(&models_da.Markup{}).Where("id = ?", id).Take(&gotMarkUp).Error)
	err := anotattionService.AddAnottation(&markUp)
	suite.Require().NoError(err)
	suite.Assert().NoError(suite.db.Model(&models_da.Markup{}).Where("id = ?", id).Take(&gotMarkUp).Error)
	suite.Assert().Equal(models_da.FromDaMarkup(&gotMarkUp), markUp)
}

func (suite *UsecaseRepositoryTestSuite) TestUsecaseDeleteMarkUp() {
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

func (suite *UsecaseRepositoryTestSuite) TestUsecaseGetMarkUp() {
	anotattionRepo := annot_repo_adapter.NewAnotattionRepositoryAdapter(suite.db)
	anotattionService := annot_service.NewAnnotattionService(anotattionRepo)
	id := uint64(1)
	markUpDa := models_da.Markup{ID: id, PageData: createPNGBuffer(TEST_VALID_PNG_IMG)}
	suite.Require().NoError(suite.db.Create(&markUpDa).Error)

	markUp, err := anotattionService.GetAnottationByID(id)
	suite.Require().NoError(err)
	suite.Require().Equal(*markUp, models_da.FromDaMarkup(&markUpDa))
}

func (suite *UsecaseRepositoryTestSuite) TestUsecaseAddMarkUpType() {
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

func (suite *UsecaseRepositoryTestSuite) TestUsecaseGetMarkUpType() {
	anotattionTypeRepo := repo_adapter.NewAnotattionTypeRepositoryAdapter(suite.db)
	anotattionTypeService := service.NewAnotattionTypeService(anotattionTypeRepo)
	id := uint64(1)
	markUpTypeDa := models_da.MarkupType{ID: id, CreatorID: 12}
	suite.Require().NoError(suite.db.Create(&markUpTypeDa).Error)

	markUpType, err := anotattionTypeService.GetAnottationTypeByID(id)
	suite.Require().NoError(err)
	suite.Require().Equal(*markUpType, models_da.FromDaMarkupType(&markUpTypeDa))
}

func (suite *UsecaseRepositoryTestSuite) TestUsecaseDeleteMarkUpType() {
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

func TestSuite(t *testing.T) {
	suite.Run(t, new(UsecaseRepositoryTestSuite))
}
