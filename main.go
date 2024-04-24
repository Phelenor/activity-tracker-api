package main

import (
	"activity-tracker-api/controllers"
	"activity-tracker-api/database"
	"activity-tracker-api/storage"
	"fmt"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"os"
	"time"
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

	app.Post("/api/login", authController.LoginHandler)

	app.Post("/api/token-refresh", authController.TokenRefreshHandler)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
	}))

	app.Use(func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims, ok := user.Claims.(jwt.MapClaims)
		if !ok {
			return c.Next()
		}

		expiry := claims["exp"].(float64)
		expiresIn := int32(time.Unix(int64(expiry), 0).Sub(time.Now()).Hours())
		if expiresIn < 24 {
			c.Set("X-Token-Expiry", fmt.Sprintf("Expires in %d hours", expiresIn))
		}

		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Idemou")
	})

	if err := app.Listen(":" + os.Getenv("API_PORT")); err != nil {
		log.Fatal(err)
	}
}
