package main

import (
	"activity-tracker-api/controllers"
	"activity-tracker-api/database"
	"activity-tracker-api/storage"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("eu-central-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), "")),
	)

	if err != nil {
		log.Fatalf("Unable to load S3 SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	s3PresignClient := s3.NewPresignClient(s3Client)
	db := database.ConnectPostgresDb()
	userRepository := storage.NewUserRepository(db)
	activityRepository := storage.NewActivityRepository(db)

	startFiberServer(userRepository, activityRepository, s3Client, s3PresignClient)
}

func startFiberServer(userRepository storage.UserRepository, activityRepository storage.ActivityRepository, s3Client *s3.Client, s3PresignClient *s3.PresignClient) {
	app := fiber.New()

	authController := controllers.AuthController{UserRepo: userRepository}
	userController := controllers.UserController{UserRepo: userRepository}
	activityController := controllers.ActivityController{ActivityRepo: activityRepository, S3Client: s3Client, S3PresignClient: s3PresignClient}

	app.Use(logger.New())

	app.Post("/api/login", authController.LoginHandler)
	app.Post("/api/token-refresh", authController.TokenRefreshHandler)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
	}))

	app.Post("/api/update-user", userController.UpdateUserDataHandler)
	app.Post("/api/delete-account", userController.DeleteAccountHandler)
	app.Post("/api/activities", activityController.PostActivityHandler)
	app.Get("/api/activities", activityController.GetActivitiesHandler)
	app.Get("/api/activities/:id", activityController.GetActivityHandler)
	app.Delete("/api/activities/:id", activityController.DeleteActivityHandler)

	if err := app.Listen(":" + os.Getenv("API_PORT")); err != nil {
		log.Fatal(err)
	}
}
