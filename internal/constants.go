package internal

import "github.com/golangbb/golangbb/v2/pkg/helpers"

var (
	PORT = helpers.GetEnv("PORT", "3000")
	DATABASE_NAME = helpers.GetEnv("DATABASE_NAME", "golangbb.db")
)
