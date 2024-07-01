package config

import (
	"fmt"
	"os"
)

var (
	DatabaseHost     = os.Getenv("DATABASE_HOST")
	DatabasePort     = os.Getenv("DATABASE_PORT")
	DatabaseUsername = os.Getenv("DATABASE_USER")
	DatabasePassword = os.Getenv("DATABASE_PASS")
	DatabaseName     = os.Getenv("DATABASE_NAME")
)

func GetDBDSN() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DatabaseUsername, DatabasePassword, DatabaseHost, DatabasePort, DatabaseName)
	return dsn
}
