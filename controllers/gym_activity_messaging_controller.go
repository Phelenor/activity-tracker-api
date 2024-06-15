package controllers

import (
	"activity-tracker-api/gym_simulator"
	"activity-tracker-api/models"
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
	gymRepo          storage.GymEquipmentRepository
	userRepo         storage.UserRepository
	connectionsMutex sync.Mutex
	connections      map[string]*websocket.Conn
	simulators       map[string]*gym_simulator.GymEquipmentSimulator
	simulatorMutex   sync.Mutex
}

func NewGymWebSocketController(gymRepo storage.GymEquipmentRepository, userRepo storage.UserRepository) *GymActivityWebSocketController {
	return &GymActivityWebSocketController{
		gymRepo:     gymRepo,
		userRepo:    userRepo,
		connections: make(map[string]*websocket.Conn),
		simulators:  make(map[string]*gym_simulator.GymEquipmentSimulator),
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

	equipment, err := controller.gymRepo.GetById(equipmentId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Send(nil)
	}

	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("userId", userId)
		c.Locals("equipmentId", strings.Clone(equipmentId))
		c.Locals("gymId", equipment.OwnerId)
		c.Locals("isGym", false)
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}

func (controller *GymActivityWebSocketController) WebSocketUpgradeHandlerUnauthorized(c *fiber.Ctx) error {
	gymId := c.Params("id")

	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("userId", "")
		c.Locals("equipmentId", "")
		c.Locals("gymId", gymId)
		c.Locals("isGym", true)
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}

func (controller *GymActivityWebSocketController) WebSocketMessageHandler(conn *websocket.Conn) {
	userId := conn.Locals("userId").(string)
	equipmentId := conn.Locals("equipmentId").(string)
	gymId := conn.Locals("gymId").(string)
	isGym := conn.Locals("isGym").(bool)

	id := equipmentId
	if isGym {
		id = gymId
	}

	controller.connectionsMutex.Lock()
	controller.connections[id] = conn
	controller.connectionsMutex.Unlock()

	if !isGym {
		controller.simulatorMutex.Lock()
		controller.simulators[equipmentId] = gym_simulator.NewGymEquipmentSimulator()
		controller.simulatorMutex.Unlock()

		user, err := controller.userRepo.GetByID(userId)
		if err != nil {
			return
		}

		equipment, err := controller.gymRepo.GetById(equipmentId)
		if err != nil {
			return
		}

		go controller.runSimulator(equipmentId, equipment.Name, gymId, *user)
	}

	defer func(isGym bool) {
		controller.connectionsMutex.Lock()
		delete(controller.connections, id)
		controller.connectionsMutex.Unlock()

		if !isGym {
			controller.simulatorMutex.Lock()
			delete(controller.simulators, equipmentId)
			controller.simulatorMutex.Unlock()
		}

		err := conn.Close()
		if err != nil {
			return
		}
	}(isGym)

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

func (controller *GymActivityWebSocketController) runSimulator(equipmentId, equipmentName, gymId string, user models.User) {
	controller.simulatorMutex.Lock()
	simulator, exists := controller.simulators[equipmentId]
	controller.simulatorMutex.Unlock()

	if !exists {
		return
	}

	for {
		if simulator.IsFinished() {
			var FinishSignal struct {
				EquipmentId string `json:"equipmentId"`
			}

			FinishSignal.EquipmentId = equipmentId
			msg, _ := json.Marshal(FinishSignal)

			if conn, ok := controller.connections[gymId]; ok {
				err := conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Error("Error sending finish signal to gym dashboard: ", err)
				}
			}

			break
		}

		if !simulator.IsActive() {
			time.Sleep(1 * time.Second)
			continue
		}

		dataSnapshot := simulator.GenerateDataSnapshot()
		dataSnapshot.UserName = user.DisplayName
		dataSnapshot.UserImageUrl = user.ImageUrl
		dataSnapshot.EquipmentId = equipmentId
		dataSnapshot.EquipmentName = equipmentName

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
