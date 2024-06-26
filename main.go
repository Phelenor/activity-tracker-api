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
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file.")
	}

	db := database.ConnectPostgresDb()
	redisStorage := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	userRepository := storage.NewUserRepository(db)
	activityRepository := storage.NewActivityRepository(db)
	groupActivityRepository := storage.NewGroupActivityRepository(db, redisStorage)
	gymAccountRepository := storage.NewGymAccountRepository(db)
	gymEquipmentRepository := storage.NewGymEquipmentRepository(db)

	s3Client, s3PresignClient := initS3()

	startPeriodicActivityClear(groupActivityRepository)

	startFiberServer(
		userRepository,
		activityRepository,
		groupActivityRepository,
		gymAccountRepository,
		gymEquipmentRepository,
		s3Client,
		s3PresignClient,
	)
}

func startFiberServer(
	userRepository storage.UserRepository,
	activityRepository storage.ActivityRepository,
	groupActivityRepository storage.GroupActivityRepository,
	gymAccountRepository storage.GymAccountRepository,
	gymEquipmentRepository storage.GymEquipmentRepository,
	s3Client *s3.Client,
	s3PresignClient *s3.PresignClient,
) {
	app := fiber.New()

	authController := controllers.NewAuthController(userRepository)
	userController := controllers.NewUserController(userRepository)
	activityController := controllers.NewActivityController(activityRepository, s3Client, s3PresignClient)
	groupActivityController := controllers.NewGroupActivityController(groupActivityRepository, userRepository)
	activityWebSocketController := controllers.NewWebSocketController(groupActivityRepository)
	gymController := controllers.NewGymController(gymAccountRepository, gymEquipmentRepository)
	gymWebSocketController := controllers.NewGymWebSocketController(gymEquipmentRepository, userRepository)

	app.Use(logger.New())
	app.Use(cors.New())

	app.Post("/api/login", authController.LoginHandler)
	app.Post("/api/token-refresh", authController.TokenRefreshHandler)

	app.Post("/api/gym/register", gymController.RegisterHandler)
	app.Post("/api/gym/login", gymController.LoginHandler)

	app.Get("/ws/activity/gym-dashboard/:id", gymWebSocketController.WebSocketUpgradeHandlerUnauthorized, websocket.New(gymWebSocketController.WebSocketMessageHandler))

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
	}))

	app.Post("/api/update-user", userController.UpdateUserDataHandler)
	app.Post("/api/delete-account", userController.DeleteAccountHandler)

	app.Post("/api/activities", activityController.PostActivityHandler)
	app.Get("/api/activities", activityController.GetActivitiesHandler)
	app.Get("/api/activities/:id", activityController.GetActivityHandler)
	app.Delete("/api/activities/:id", activityController.DeleteActivityHandler)

	app.Post("/api/create-group-activity", groupActivityController.CreateGroupActivityHandler)
	app.Post("/api/join-group-activity", groupActivityController.JoinGroupActivityHandler)
	app.Post("/api/leave-group-activity", groupActivityController.LeaveGroupActivityHandler)
	app.Get("/api/group-activities/:id", groupActivityController.GetGroupActivityHandler)
	app.Delete("/api/group-activities/:id", groupActivityController.DeleteGroupActivityHandler)
	app.Get("/api/group-activities", groupActivityController.GetScheduledActivitiesHandler)
	app.Get("/api/group-activity-overview/:id", groupActivityController.GetGroupActivityOverviewHandler)
	app.Get("/ws/activity/:id", activityWebSocketController.WebSocketUpgradeHandler, websocket.New(activityWebSocketController.WebSocketMessageHandler))

	app.Get("/api/gym/equipment", gymController.GetAllEquipmentHandler)
	app.Post("/api/gym/equipment", gymController.CreateEquipmentHandler)
	app.Get("/api/gym/equipment/:id", gymController.GetEquipmentHandler)
	app.Get("/ws/activity/gym/:id", gymWebSocketController.WebSocketUpgradeHandler, websocket.New(gymWebSocketController.WebSocketMessageHandler))

	if err := app.Listen(":" + os.Getenv("API_PORT")); err != nil {
		log.Fatal(err)
	}
}

func initS3() (*s3.Client, *s3.PresignClient) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("eu-central-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), ""),
		),
	)

	if err != nil {
		log.Fatalf("Unable to load S3 SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	s3PresignClient := s3.NewPresignClient(s3Client)

	return s3Client, s3PresignClient
}

func startPeriodicActivityClear(repo storage.GroupActivityRepository) {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			err := repo.DeleteExpiredActivities()
			if err != nil {
				log.Error(err)
			}
		}
	}()
}
