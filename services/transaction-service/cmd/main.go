package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rezkyauliapratama/fsi-playground/services/transaction-service/config"
	"github.com/rezkyauliapratama/fsi-playground/services/transaction-service/database"
	"github.com/rezkyauliapratama/fsi-playground/services/transaction-service/handlers"
	"github.com/rezkyauliapratama/fsi-playground/services/transaction-service/repositories"
	"github.com/rezkyauliapratama/fsi-playground/services/transaction-service/services"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())
	db := database.InitDB(config.GetDBDSN())
	defer db.Close()
	accountRepo := repositories.NewAccountRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	entryRepo := repositories.NewEntryRepository(db)
	transactionService := services.NewTransactionService(transactionRepo, accountRepo, entryRepo)
	handlers.RegisterTransactionHandlers(app, transactionService)

	log.Fatal(app.Listen(":8000"))

}
