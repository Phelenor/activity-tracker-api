package controllers

import (
	"activity-tracker-api/gymSimulator"
	"activity-tracker-api/models/ws"
	"activity-tracker-api/storage"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"sync"
	"time"
)

type GymActivityWebSocketController struct {
	repository       storage.GymEquipmentRepository
	connectionsMutex sync.Mutex
	connections      map[string]*websocket.Conn
	simulators       map[string]*gymSimulator.GymEquipmentSimulator
	simulatorMutex   sync.Mutex
}

func NewGymWebSocketController(repository storage.GymEquipmentRepository) *GymActivityWebSocketController {
	return &GymActivityWebSocketController{
		repository:  repository,
		connections: make(map[string]*websocket.Conn),
		simulators:  make(map[string]*gymSimulator.GymEquipmentSimulator),
	}
}

func (controller *GymActivityWebSocketController) WebSocketUpgradeHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).Send(nil)
	}

	userId := claims["id"].(string)
	equipmentId := c.Params("id")

	equipment, err := controller.repository.GetById(equipmentId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Send(nil)
	}

	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("userId", userId)
		c.Locals("equipmentId", strings.Clone(equipmentId))
		c.Locals("gymId", equipment.OwnerId)
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}

func (controller *GymActivityWebSocketController) WebSocketMessageHandler(conn *websocket.Conn) {
	equipmentId := conn.Locals("equipmentId").(string)
	gymId := conn.Locals("gymId").(string)

	controller.connectionsMutex.Lock()
	controller.connections[equipmentId] = conn
	controller.connectionsMutex.Unlock()

	controller.simulatorMutex.Lock()
	controller.simulators[equipmentId] = gymSimulator.NewGymEquipmentSimulator()
	controller.simulatorMutex.Unlock()

	go controller.runSimulator(equipmentId, gymId)

	defer func() {
		controller.connectionsMutex.Lock()
		delete(controller.connections, equipmentId)
		controller.connectionsMutex.Unlock()

		controller.simulatorMutex.Lock()
		delete(controller.simulators, equipmentId)
		controller.simulatorMutex.Unlock()

		err := conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		controller.handleIncomingMessage(equipmentId, msg)
	}
}

func (controller *GymActivityWebSocketController) handleIncomingMessage(equipmentId string, msg []byte) {
	var msgType struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(msg, &msgType); err != nil {
		log.Error("No type provided in the message: ", err)
		return
	}

	switch msgType.Type {
	case "control_action":
		var controlAction ws.ControlAction
		if err := json.Unmarshal(msg, &controlAction); err != nil {
			log.Error("Error unmarshalling ControlAction: ", err)
			return
		}

		controller.simulatorMutex.Lock()
		simulator, exists := controller.simulators[equipmentId]
		controller.simulatorMutex.Unlock()

		if !exists {
			log.Error("Simulator not found for equipmentId: ", equipmentId)
			return
		}

		switch controlAction.Action {
		case ws.ActivityControlStart:
			simulator.Start()
		case ws.ActivityControlPause:
			simulator.Pause()
		case ws.ActivityControlResume:
			simulator.Resume()
		case ws.ActivityControlFinish:
			simulator.Finish()
		default:
			log.Error("Unknown control action: ", controlAction.Action)
		}

	default:
		log.Error("Unknown message type: ", msgType.Type)
	}
}

func (controller *GymActivityWebSocketController) runSimulator(equipmentId, gymId string) {
	controller.simulatorMutex.Lock()
	simulator, exists := controller.simulators[equipmentId]
	controller.simulatorMutex.Unlock()

	if !exists {
		return
	}

	for {
		if !simulator.IsActive() {
			time.Sleep(1 * time.Second)
			continue
		}

		if simulator.IsFinished() {
			break
		}

		dataSnapshot := simulator.GenerateDataSnapshot()

		msg, err := json.Marshal(dataSnapshot)
		if err != nil {
			log.Error("Error marshaling DataSnapshot: ", err)
			continue
		}

		controller.connectionsMutex.Lock()
		if conn, ok := controller.connections[equipmentId]; ok {
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Error("Error sending data to client app: ", err)
			}
		}

		if conn, ok := controller.connections[gymId]; ok {
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Error("Error sending data to gym dashboard: ", err)
			}
		}

		controller.connectionsMutex.Unlock()

		time.Sleep(1 * time.Second)
	}
}
