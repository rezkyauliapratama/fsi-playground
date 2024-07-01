package config

import "os"

var (
	DatabaseHost     = os.Getenv("DATABASE_HOST")
	DatabasePort     = os.Getenv("DATABASE_PORT")
	DatabaseUsername = os.Getenv("DATABASE_USER")
	DatabasePassword = os.Getenv("DATABASE_PASS")
	DatabaseName     = os.Getenv("DATABASE_NAME")
)
