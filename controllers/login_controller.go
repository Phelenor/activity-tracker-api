package controllers

import (
	"activity-tracker-api/models"
	"activity-tracker-api/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserController struct {
	UserRepo storage.UserRepository
}

func (controller *UserController) ChangeNameHandler(c *fiber.Ctx) error {
	request := models.UserChangeNameRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request.")
	}

	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Next()
	}

	userId := claims["id"].(string)
	dbUser, err := controller.UserRepo.GetByID(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	if dbUser == nil {
		return c.Status(fiber.StatusNotFound).SendString("User not found.")
	}

	dbUser.DisplayName = request.Name

	if err := controller.UserRepo.Update(dbUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	return c.Status(fiber.StatusOK).JSON(dbUser)
}

func (controller *UserController) DeleteAccountHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Next()
	}

	userId := claims["id"].(string)

	if err := controller.UserRepo.Delete(userId); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	return c.Status(fiber.StatusOK).SendString("User deleted.")
}
