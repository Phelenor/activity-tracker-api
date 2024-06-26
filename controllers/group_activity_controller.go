package controllers

import (
	"activity-tracker-api/models/activity"
	"activity-tracker-api/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"math/rand"
	"slices"
	"sync"
)

type GroupActivityController struct {
	GroupActivityRepo storage.GroupActivityRepository
	UserRepo          storage.UserRepository
}

func NewGroupActivityController(groupActivityRepo storage.GroupActivityRepository, userRepo storage.UserRepository) *GroupActivityController {
	return &GroupActivityController{
		GroupActivityRepo: groupActivityRepo,
		UserRepo:          userRepo,
	}
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
		ConnectedUsers: []string{},
		FinishedUsers:  []string{},
		ActiveUsers:    []string{},
	}

	err = controller.GroupActivityRepo.InsertIntoRedis(&groupActivity)
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

	err = controller.GroupActivityRepo.InsertIntoRedis(groupActivity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(".")
	}

	return c.Status(fiber.StatusOK).JSON(groupActivity)
}

func (controller *GroupActivityController) LeaveGroupActivityHandler(c *fiber.Ctx) error {
	request := activity.LeaveGroupActivityRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request.")
	}

	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).Send(nil)
	}

	userId := claims["id"].(string)

	err := controller.GroupActivityRepo.RemoveUserFromActivityList(request.ActivityId, userId, storage.ActivityListTypeJoined)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).Send(nil)
}

func (controller *GroupActivityController) GetGroupActivityHandler(c *fiber.Ctx) error {
	activityId := c.Params("id")

	groupActivity, err := controller.GroupActivityRepo.GetByIDFromRedis(activityId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Send(nil)
	}

	return c.Status(fiber.StatusOK).JSON(groupActivity)
}

func (controller *GroupActivityController) GetGroupActivityOverviewHandler(c *fiber.Ctx) error {
	activityId := c.Params("id")

	groupActivity, err := controller.GroupActivityRepo.GetByIDFromDb(activityId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Send(nil)
	}

	groupActivityOverview := activity.GroupActivityOverview{OwnerId: groupActivity.UserOwnerId}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, userId := range groupActivity.FinishedUsers {
		wg.Add(1)
		go func(userId string) {
			defer wg.Done()
			user, err := controller.UserRepo.GetByID(userId)
			if err == nil {
				mutex.Lock()
				groupActivityOverview.Users = append(groupActivityOverview.Users, *user)
				mutex.Unlock()
			}
		}(userId)
	}

	wg.Wait()

	return c.Status(fiber.StatusOK).JSON(groupActivityOverview)
}

func (controller *GroupActivityController) DeleteGroupActivityHandler(c *fiber.Ctx) error {
	activityId := c.Params("id")
	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).Send(nil)
	}

	userId := claims["id"].(string)

	groupActivity, err := controller.GroupActivityRepo.GetByIDFromRedis(activityId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Send(nil)
	}

	if groupActivity.UserOwnerId != userId {
		return c.Status(fiber.StatusForbidden).Send(nil)
	}

	err = controller.GroupActivityRepo.Delete(activityId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Send(nil)
	}

	return c.Status(fiber.StatusOK).Send(nil)
}

func (controller *GroupActivityController) GetScheduledActivitiesHandler(c *fiber.Ctx) error {
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

	if len(groupActivities) == 0 {
		return c.Status(fiber.StatusOK).JSON(make([]*string, 0))
	} else {
		return c.Status(fiber.StatusOK).JSON(groupActivities)
	}
}

func generateJoinCode(length int) string {
	chars := "0123456789ABCDEF"
	code := make([]byte, length)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}

	return string(code)
}
