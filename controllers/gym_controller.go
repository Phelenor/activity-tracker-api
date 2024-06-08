package controllers

import (
	"activity-tracker-api/models"
	"activity-tracker-api/storage"
	"activity-tracker-api/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GymController struct {
	GymRepository storage.GymRepository
}

func NewGymController(repository storage.GymRepository) *GymController {
	return &GymController{
		GymRepository: repository,
	}
}

func (controller *GymController) RegisterHandler(c *fiber.Ctx) error {
	request := models.GymRegisterRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request.")
	}

	account := models.GymAccount{
		Id:           uuid.New().String(),
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: request.PasswordHash,
	}

	err := controller.GymRepository.Insert(&account)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	return c.Status(fiber.StatusOK).Send(nil)
}

func (controller *GymController) LoginHandler(c *fiber.Ctx) error {
	request := models.GymLoginRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request.")
	}

	accountDb, err := controller.GymRepository.GetByEmail(request.Email)
	if err != nil || accountDb.PasswordHash != request.PasswordHash {
		return c.Status(fiber.StatusUnauthorized).Send(nil)
	}

	accessToken, err := util.BuildAccessToken(accountDb.Id, accountDb.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	response := models.GymTokenResponse{
		GymAccount:  *accountDb,
		AccessToken: accessToken,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
