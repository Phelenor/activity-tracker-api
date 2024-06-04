package controllers

import (
	"activity-tracker-api/models/activity"
	"activity-tracker-api/models/ws"
	"activity-tracker-api/storage"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"sync"
)

type ActivityWebSocketController struct {
	GroupActivityRepo storage.GroupActivityRepository
	connectionsMutex  sync.Mutex
	connections       map[string]*websocket.Conn
}

func NewWebSocketController(repository storage.GroupActivityRepository) *ActivityWebSocketController {
	return &ActivityWebSocketController{
		GroupActivityRepo: repository,
		connections:       make(map[string]*websocket.Conn),
	}
}

func (controller *ActivityWebSocketController) WebSocketUpgradeHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).Send(nil)
	}

	userId := claims["id"].(string)
	activityId := c.Params("id")

	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("userId", userId)
		c.Locals("activityId", strings.Clone(activityId))
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}

func (controller *ActivityWebSocketController) WebSocketMessageHandler(conn *websocket.Conn) {
	userId := conn.Locals("userId").(string)
	activityId := conn.Locals("activityId").(string)

	controller.connectionsMutex.Lock()
	controller.connections[userId] = conn
	controller.connectionsMutex.Unlock()
	err := controller.GroupActivityRepo.AddUserToActivityList(activityId, userId, storage.ActivityListTypeConnected)
	if err != nil {
		log.Error("error connecting user to activity: ", err)
	}

	controller.broadcastActivityUpdate(activityId)

	defer func(activityId, userId string) {
		controller.connectionsMutex.Lock()
		delete(controller.connections, userId)
		controller.connectionsMutex.Unlock()
		err := controller.GroupActivityRepo.RemoveUserFromActivityList(activityId, userId, storage.ActivityListTypeConnected)
		if err != nil {
			log.Error("error disconnecting user from activity: ", err)
		}

		controller.broadcastActivityUpdate(activityId)

		err = conn.Close()
		if err != nil {
			log.Error(err)
		}
	}(activityId, userId)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		controller.handleIncomingMessage(activityId, userId, msg)
	}
}

func (controller *ActivityWebSocketController) handleIncomingMessage(activityId string, userId string, msg []byte) {
	var msgType struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(msg, &msgType); err != nil {
		log.Error("No type provided in the message: ", err)
		return
	}

	switch msgType.Type {
	case "user_update":
		controller.broadcastMessage(activityId, userId, msg)
	case "control_action":
		controller.broadcastMessage(activityId, userId, msg)

		var controlAction ws.ControlAction
		if err := json.Unmarshal(msg, &controlAction); err != nil {
			log.Error("Error unmarshalling ControlAction: ", err)
			return
		}

		var status activity.ActivityStatus

		switch controlAction.Action {
		case ws.ActivityControlResume:
		case ws.ActivityControlStart:
			status = activity.ActivityStatusInProgress
		case ws.ActivityControlPause:
			status = activity.ActivityStatusPaused
		case ws.ActivityControlFinish:
			status = activity.ActivityStatusFinished
		default:
			status = activity.ActivityStatusUndefined
		}

		err := controller.GroupActivityRepo.UpdateActivityStatus(activityId, status)
		if err != nil {
			log.Error("Error updating activity status: ", err)
		}
	case "user_finish_signal":
		controller.broadcastMessage(activityId, userId, msg)

		var userFinish ws.UserFinish
		if err := json.Unmarshal(msg, &userFinish); err != nil {
			log.Error("Error unmarshalling UserFinish: ", err)
			return
		}

		err := controller.GroupActivityRepo.RemoveUserFromActivityList(activityId, userId, storage.ActivityListTypeActive)
		if err != nil {
			return
		}
	default:
		log.Error("Unknown message type: ", msgType.Type)
	}
}

func (controller *ActivityWebSocketController) broadcastMessage(activityId string, senderUserId string, msg []byte) {
	groupActivity, err := controller.GroupActivityRepo.GetByIDFromRedis(activityId)
	if err != nil {
		return
	}

	for _, userId := range groupActivity.ConnectedUsers {
		if userId == senderUserId {
			continue
		}

		controller.connectionsMutex.Lock()
		conn, ok := controller.connections[userId]
		controller.connectionsMutex.Unlock()

		if ok {
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (controller *ActivityWebSocketController) broadcastActivityUpdate(activityId string) {
	groupActivity, err := controller.GroupActivityRepo.GetByIDFromRedis(activityId)
	if err != nil {
		return
	}

	msg := ws.ActivityUpdate{
		Activity: *groupActivity,
		Type:     "activity_update",
	}

	msgJson, err := json.Marshal(msg)
	if err != nil {
		log.Error("Error marshalling ActivityUpdate: ", err)
		return
	}

	for _, userId := range groupActivity.ConnectedUsers {
		controller.connectionsMutex.Lock()
		conn, ok := controller.connections[userId]
		controller.connectionsMutex.Unlock()

		if ok {
			err := conn.WriteMessage(websocket.TextMessage, msgJson)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
