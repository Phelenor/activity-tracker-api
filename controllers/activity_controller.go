package controllers

import (
	"activity-tracker-api/models"
	"activity-tracker-api/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ActivityController struct {
	UserRepo storage.UserRepository
}

func (controller *ActivityController) PostActivityHandler(c *fiber.Ctx) error {
	request := models.Activity{}

	if err := c.BodyParser(&request); err != nil {
		log.Debug(err.Error())
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request.")
	}

	log.Debug(request)

	return c.Status(fiber.StatusOK).JSON(request)
}
