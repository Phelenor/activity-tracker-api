package middleware

import (
	"activity-tracker-api/models"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func CreateUserJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return signed, err
}
