package app

import (
	"activity-tracker-api/controllers"
	"activity-tracker-api/database"
	"activity-tracker-api/storage"
	jwtware "github.com/gofiber/contrib/jwt"
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

	app.Get("/unrestricted", func(c *fiber.Ctx) error {
		return c.SendString("hello world unrestricted")
	})

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
	}))

	app.Get("/restricted", func(c *fiber.Ctx) error {
		return c.SendString("hello world restricted")
	})

	if err := app.Listen(":" + os.Getenv("API_PORT")); err != nil {
		log.Fatal(err)
	}
}
