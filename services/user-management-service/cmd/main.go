package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rezkyauliapratama/fsi-playground/services/user-management-service/config"
	"github.com/rezkyauliapratama/fsi-playground/services/user-management-service/database"
	"github.com/rezkyauliapratama/fsi-playground/services/user-management-service/handlers"
	"github.com/rezkyauliapratama/fsi-playground/services/user-management-service/repositories"
	"github.com/rezkyauliapratama/fsi-playground/services/user-management-service/services"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())
	db := database.InitDB(config.GetDBDSN())
	defer db.Close()
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	handlers.RegisterUserHandlers(app, userService)

	log.Fatal(app.Listen(":8001"))

}
