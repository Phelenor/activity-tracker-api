package controllers

import (
	"activity-tracker-api/models"
	"activity-tracker-api/storage"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/api/idtoken"
	"os"
)

func SetupAuthController(app *fiber.App, userRepo storage.UserRepository) {
	app.Post("/login", func(c *fiber.Ctx) error {
		return login(c, userRepo)
	})
}

func login(c *fiber.Ctx, userRepo storage.UserRepository) error {
	request := models.LoginRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request.")
	}

	payload, err := idtoken.Validate(context.Background(), request.IdToken, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Can't validate token.")
	}

	if payload.Claims["nonce"] != request.Nonce {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid token nonce.")
	}

	user := models.User{
		Id:       payload.Claims["sub"].(string),
		Name:     payload.Claims["name"].(string),
		Email:    payload.Claims["email"].(string),
		ImageUrl: payload.Claims["picture"].(string),
	}

	dbUser, err := userRepo.GetByID(user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	if dbUser == nil {
		if err := userRepo.Insert(&user); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
		}

		log.Info("User registered. ", user.Email)
	} else {
		if err := userRepo.Update(&user); err != nil {
			log.Error("User update error. ", err)
		}

		log.Info("User login. ", user.Email)
	}

	return nil
}
