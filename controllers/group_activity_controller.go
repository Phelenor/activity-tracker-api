package controllers

import (
	"activity-tracker-api/models/activity"
	"activity-tracker-api/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"math/rand"
	"slices"
)

type GroupActivityController struct {
	GroupActivityRepo storage.GroupActivityRepository
	UserRepo          storage.UserRepository
}

func (controller *GroupActivityController) CreateGroupActivityHandler(c *fiber.Ctx) error {
	request := activity.CreateGroupActivityRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request.")
	}

	userToken := c.Locals("user").(*jwt.Token)
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).Send(nil)
	}

	userId := claims["id"].(string)

	user, err := controller.UserRepo.GetByID(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error while saving activity.")
	}

	groupActivity := activity.GroupActivity{
		Id:             uuid.New().String(),
		JoinCode:       generateJoinCode(6),
		UserOwnerId:    userId,
		UserOwnerName:  user.DisplayName,
		ActivityType:   request.ActivityType,
		StartTimestamp: request.StartTimestamp,
		Status:         activity.ActivityStatusNotStarted,
		JoinedUsers:    []string{userId},
		StartedUsers:   []string{},
		ActiveUsers:    []string{},
	}

	err = controller.GroupActivityRepo.Insert(&groupActivity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error while saving activity.")
	}

	return c.Status(fiber.StatusOK).JSON(groupActivity)
}

func (controller *GroupActivityController) JoinGroupActivityHandler(c *fiber.Ctx) error {
	request := activity.JoinGroupActivityRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request.")
	}

	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).Send(nil)
	}

	userId := claims["id"].(string)

	groupActivity, err := controller.GroupActivityRepo.GetByJoinCodeFromRedis(request.JoinCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Invalid join code.")
	}

	if !slices.Contains(groupActivity.JoinedUsers, userId) {
		groupActivity.JoinedUsers = append(groupActivity.JoinedUsers, userId)
	}

	err = controller.GroupActivityRepo.Insert(groupActivity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(".")
	}

	return c.Status(fiber.StatusOK).JSON(groupActivity)
}

func (controller *GroupActivityController) GetGroupActivityHandler(c *fiber.Ctx) error {
	activityId := c.Params("id")

	groupActivity, err := controller.GroupActivityRepo.GetByIDFromRedis(activityId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Send(nil)
	}

	return c.Status(fiber.StatusOK).JSON(groupActivity)
}

func (controller *GroupActivityController) GetPendingActivitiesHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).Send(nil)
	}

	userId := claims["id"].(string)

	groupActivities, err := controller.GroupActivityRepo.GetByUserIdFromRedis(userId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Send(nil)
	}

	return c.Status(fiber.StatusOK).JSON(groupActivities)
}

func generateJoinCode(length int) string {
	chars := "0123456789ABCDEF"
	code := make([]byte, length)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}

	return string(code)
}
