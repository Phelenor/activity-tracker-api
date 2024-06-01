package controllers

import (
	"activity-tracker-api/models/ws"
	"activity-tracker-api/storage"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
	"sync"
)

type ActivityWebSocketController struct {
	GroupActivityRepo storage.GroupActivityRepository
	connections       map[string]*websocket.Conn
	connectionsMutex  sync.Mutex
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

	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("userId", userId)
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}

func (controller *ActivityWebSocketController) WebSocketMessageHandler(c *websocket.Conn) {
	userId := c.Locals("userId").(string)

	controller.connectionsMutex.Lock()
	controller.connections[userId] = c
	controller.connectionsMutex.Unlock()

	defer func() {
		controller.connectionsMutex.Lock()
		delete(controller.connections, userId)
		controller.connectionsMutex.Unlock()
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break
		}

		controller.handleIncomingMessage(msg)
	}
}

func (controller *ActivityWebSocketController) handleIncomingMessage(msg []byte) {
	log.Debug("new message: ", string(msg))

	var baseMsg ws.ActivityMessage
	if err := json.Unmarshal(msg, &baseMsg); err != nil {
		log.Debug("Error unmarshaling base message: %v\n", err)
		return
	}
}
