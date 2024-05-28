package controllers

import (
	"activity-tracker-api/models/activity"
	"activity-tracker-api/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"math/rand"
)

type GroupActivityController struct {
	GroupActivityRepo storage.GroupActivityRepository
}

func (controller *GroupActivityController) CreateGroupActivityHandler(c *fiber.Ctx) error {
	request := activity.CreateGroupActivityRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request.")
	}

	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Next()
	}

	userId := claims["id"].(string)

	groupActivity := activity.GroupActivity{
		Id:             uuid.New().String(),
		JoinCode:       generateJoinCode(6),
		UserOwnerId:    userId,
		ActivityType:   request.ActivityType,
		StartTimestamp: request.StartTimestamp,
		Status:         activity.ActivityStatusNotStarted,
		StartedUsers:   []string{userId},
		ActiveUsers:    []string{userId},
	}

	err := controller.GroupActivityRepo.Insert(&groupActivity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error while saving activity.")
	}

	return c.Status(fiber.StatusOK).JSON(groupActivity)
}

func generateJoinCode(length int) string {
	digits := "0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
	}

	return string(code)
}
