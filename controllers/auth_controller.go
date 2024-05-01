package controllers

import (
	"activity-tracker-api/models"
	"activity-tracker-api/storage"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
	"os"
	"time"
)

func buildAccessToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return signed, err
}

func buildRefreshToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(time.Hour * 24 * 14).Unix(),
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
		return c.Status(fiber.StatusBadRequest).SendString("Can't validate accessToken.")
	}

	if payload.Claims["nonce"] != request.Nonce {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid accessToken nonce.")
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

	accessToken, err := buildAccessToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	refreshToken, err := buildRefreshToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	response := models.UserTokenResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (controller *AuthController) TokenRefreshHandler(c *fiber.Ctx) error {
	request := models.TokenRefreshRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request.")
	}

	token, err := jwt.Parse(request.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token.")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid token.")
	}

	userId := claims["id"].(string)

	user, err := controller.UserRepo.GetByID(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	if user == nil {
		return c.Status(fiber.StatusBadRequest).SendString("User not found.")
	}

	accessToken, err := buildAccessToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	refreshToken, err := buildRefreshToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	response := models.UserTokenResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
