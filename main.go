package main

import (
	"activity-tracker-api/controllers"
	"activity-tracker-api/database"
	"activity-tracker-api/storage"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file.")
	}

	db := database.ConnectPostgresDb()
	userRepository := storage.NewUserRepository(db)

	startFiberServer(userRepository)
}

func startFiberServer(userRepository storage.UserRepository) {
	app := fiber.New()

	authController := controllers.AuthController{UserRepo: userRepository}
	userController := controllers.UserController{UserRepo: userRepository}
	activityController := controllers.ActivityController{UserRepo: userRepository}

	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\nRequest:\t${body}\nResponse:\t${resBody}\n",
	}))

	app.Post("/api/login", authController.LoginHandler)
	app.Post("/api/token-refresh", authController.TokenRefreshHandler)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
	}))

	app.Post("/api/update-user", userController.UpdateUserDataHandler)
	app.Post("/api/delete-account", userController.DeleteAccountHandler)

	app.Post("/api/activity", activityController.PostActivityHandler)

	if err := app.Listen(":" + os.Getenv("API_PORT")); err != nil {
		log.Fatal(err)
	}
}
