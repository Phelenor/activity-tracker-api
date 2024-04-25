package controllers

import (
	"activity-tracker-api/models"
	"activity-tracker-api/storage"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
	"os"
	"time"
)

func buildUserJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return signed, err
}

type AuthController struct {
	UserRepo storage.UserRepository
}

func (controller *AuthController) LoginHandler(c *fiber.Ctx) error {
	request := models.LoginRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request.")
	}

	payload, err := idtoken.Validate(context.Background(), request.IdToken, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Can't validate token.")
	}

	if payload.Claims["nonce"] != request.Nonce {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid token nonce.")
	}

	user := models.User{
		Id:          payload.Claims["sub"].(string),
		Name:        payload.Claims["name"].(string),
		DisplayName: payload.Claims["name"].(string),
		Email:       payload.Claims["email"].(string),
		ImageUrl:    payload.Claims["picture"].(string),
	}

	dbUser, err := controller.UserRepo.GetByID(user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	if dbUser == nil {
		err = controller.UserRepo.Insert(&user)
	} else {
		user.DisplayName = dbUser.DisplayName
		err = controller.UserRepo.Update(&user)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	token, err := buildUserJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	response := models.UserTokenResponse{
		User:  user,
		Token: token,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (controller *AuthController) TokenRefreshHandler(c *fiber.Ctx) error {
	request := models.TokenRefreshRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request.")
	}

	user, err := controller.UserRepo.GetByID(request.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	if user == nil {
		return c.Status(fiber.StatusBadRequest).SendString("User not found.")
	}

	token, err := buildUserJWT(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	response := models.UserTokenResponse{
		User:  *user,
		Token: token,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
