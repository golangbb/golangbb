package internal

import (
	"github.com/golangbb/golangbb/v2/pkg/helpers"
)

var (
	keyPORT             = "PORT"
	defaultPORT         = "3000"
	keyDATABASENAME     = "DATABASENAME"
	defaultDATABASENAME = "golang.db"

	PORT         = helpers.GetEnv(keyPORT, defaultPORT)
	DATABASENAME = helpers.GetEnv(keyDATABASENAME, defaultDATABASENAME)
)
