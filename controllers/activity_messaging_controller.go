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
	activityId := c.Params("id")

	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("userId", userId)
		c.Locals("activityId", activityId)
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}

func (controller *ActivityWebSocketController) WebSocketMessageHandler(conn *websocket.Conn) {
	userId := conn.Locals("userId").(string)
	// activityId := conn.Locals("activityId").(string)

	controller.connectionsMutex.Lock()
	controller.connections[userId] = conn
	controller.connectionsMutex.Unlock()

	defer func() {
		controller.connectionsMutex.Lock()
		delete(controller.connections, userId)
		controller.connectionsMutex.Unlock()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		log.Debug("new message: ", string(msg))

		err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Fatal(err)
		}

		// controller.handleIncomingMessage(msg)
	}
}

func (controller *ActivityWebSocketController) handleIncomingMessage(msg []byte) {
	log.Debug("new message: ", string(msg))

	var baseMsg ws.ActivityMessage
	if err := json.Unmarshal(msg, &baseMsg); err != nil {
		log.Debug("Error unmarshaling base message: %v\n", err)
		return
	}

	//switch baseMsg.Type {
	//case "connect_message":
	//	var aMsg ConnectMessage
	//	if err := json.Unmarshal(baseMsg.Data, &aMsg); err != nil {
	//		log.Printf("Error unmarshaling ConnectMessage: %v\n", err)
	//		return
	//	}
	//	controller.handleConnectMessage(userID, aMsg)
	//case "data_update":
	//	var bMsg DataUpdate
	//	if err := json.Unmarshal(baseMsg.Data, &bMsg); err != nil {
	//		log.Printf("Error unmarshaling DataUpdate: %v\n", err)
	//		return
	//	}
	//	controller.handleDataUpdate(userID, bMsg)
	//case "status_change":
	//	var cMsg StatusChange
	//	if err := json.Unmarshal(baseMsg.Data, &cMsg); err != nil {
	//		log.Printf("Error unmarshaling StatusChange: %v\n", err)
	//		return
	//	}
	//	controller.handleStatusChange(userID, cMsg)
	//default:
	//	log.Printf("Unknown message type: %s\n", baseMsg.Type)
	//}
}
