package app

import (
	"activity-tracker-api/controllers"
	"activity-tracker-api/database"
	"activity-tracker-api/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"os"
)

func SetupAndRun() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file.")
	}

	db := database.ConnectPostgresDb()
	userRepository := storage.NewUserRepository(db)

	startFiberServer(userRepository)
}

func startFiberServer(userRepository storage.UserRepository) {
	app := fiber.New()

	controllers.SetupAuthController(app, userRepository)

	if err := app.Listen(":" + os.Getenv("API_PORT")); err != nil {
		log.Fatal(err)
	}
}
