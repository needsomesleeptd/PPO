package integration_utils

import (
	"gorm.io/driver/postgres"
)

var (
	CONN_STR = "host=localhost user=andrew password=1 database=lab01db port=5432"
	TestCfg  = postgres.Config{DSN: CONN_STR}
)
