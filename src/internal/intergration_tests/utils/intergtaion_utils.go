package integration_utils

import (
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	CONN_STR = "host=localhost user=andrew password=1 database=lab01db port=5432"
	TestCfg  = postgres.Config{DSN: CONN_STR}
)

type UsecaseRepositoryTestSuite struct {
	suite.Suite
	db *gorm.DB
}
