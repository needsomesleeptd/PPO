package integration_tests

import (
	auth_service "annotater/internal/bl/auth"
	user_repo_adapter "annotater/internal/bl/userService/userRepo/userRepoAdapter"
	integration_utils "annotater/internal/intergration_tests/utils"
	"annotater/internal/models"
	models_da "annotater/internal/models/modelsDA"
	auth_utils "annotater/internal/pkg/authUtils"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AuthTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *AuthTestSuite) SetupTest() {
	// Open a new database connection for each test
	db, err := gorm.Open(postgres.New(integration_utils.TestCfg), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// Automatically migrate the schema for each test
	db.AutoMigrate(&models_da.User{})
	suite.db = db
}

func (suite *AuthTestSuite) TearDownTest() {
	// Delete the test table after each test
	suite.db.Migrator().DropTable(&models_da.User{})
}

func (suite *AuthTestSuite) TestUsecaseSignUp() {
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

func (suite *AuthTestSuite) TestUsecaseSignIn() {
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

func TestSuiteAuth(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
