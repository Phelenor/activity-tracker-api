package controllers

import (
	"activity-tracker-api/models"
	"activity-tracker-api/models/gym"
	"activity-tracker-api/storage"
	"activity-tracker-api/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type GymController struct {
	GymAccountRepository   storage.GymAccountRepository
	GymEquipmentRepository storage.GymEquipmentRepository
}

func NewGymController(accountRepo storage.GymAccountRepository, equipmentRepo storage.GymEquipmentRepository) *GymController {
	return &GymController{
		GymAccountRepository:   accountRepo,
		GymEquipmentRepository: equipmentRepo,
	}
}

func (controller *GymController) RegisterHandler(c *fiber.Ctx) error {
	request := models.GymRegisterRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request.")
	}

	account := gym.GymAccount{
		Id:           uuid.New().String(),
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: request.PasswordHash,
	}

	err := controller.GymAccountRepository.Insert(&account)

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

	accountDb, err := controller.GymAccountRepository.GetByEmail(request.Email)
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

func (controller *GymController) GetAllEquipmentHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).Send(nil)
	}

	userId := claims["id"].(string)

	nameQuery := c.Query("q", "")

	equipment, err := controller.GymEquipmentRepository.GetForUserId(userId, nameQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	return c.Status(fiber.StatusOK).JSON(equipment)
}

func (controller *GymController) GetEquipmentHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	equipment, err := controller.GymEquipmentRepository.GetById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	return c.Status(fiber.StatusOK).JSON(equipment)
}

func (controller *GymController) CreateEquipmentHandler(c *fiber.Ctx) error {
	request := gym.CreateEquipmentRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request.")
	}

	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).Send(nil)
	}

	userId := claims["id"].(string)

	equipment := gym.GymEquipment{
		Id:           uuid.New().String(),
		OwnerId:      userId,
		Name:         request.Name,
		Description:  request.Description,
		ImageUrl:     request.ImageUrl,
		VideoUrl:     request.VideoUrl,
		ActivityType: request.ActivityType,
	}

	err := controller.GymEquipmentRepository.Insert(&equipment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	return c.Status(fiber.StatusOK).JSON(equipment)
}
