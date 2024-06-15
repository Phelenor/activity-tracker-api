package controllers

import (
	"activity-tracker-api/models"
	"activity-tracker-api/storage"
	"activity-tracker-api/util"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
	"os"
)

type AuthController struct {
	UserRepo storage.UserRepository
}

func NewAuthController(repository storage.UserRepository) *AuthController {
	return &AuthController{
		UserRepo: repository,
	}
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

	imageUrl, ok := payload.Claims["picture"].(string)
	if !ok {
		imageUrl = ""
	}

	user := models.User{
		Id:          payload.Claims["sub"].(string),
		Name:        payload.Claims["name"].(string),
		DisplayName: payload.Claims["name"].(string),
		Email:       payload.Claims["email"].(string),
		ImageUrl:    imageUrl,
	}

	dbUser, err := controller.UserRepo.GetByID(user.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	if dbUser == nil {
		err = controller.UserRepo.Insert(&user)
	} else {
		user.DisplayName = dbUser.DisplayName
		user.Weight = dbUser.Weight
		user.Height = dbUser.Height
		user.BirthTimestamp = dbUser.BirthTimestamp
		err = controller.UserRepo.Update(&user)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	accessToken, err := util.BuildAccessToken(user.Id, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	refreshToken, err := util.BuildRefreshToken(user.Id)
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

	accessToken, err := util.BuildAccessToken((*user).Id, (*user).Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error.")
	}

	refreshToken, err := util.BuildRefreshToken((*user).Id)
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
