package helpers

import (
	"gojinmongo/models"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(user *models.User) string {
	// TODO: Generate JWT token here
	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"name":    user.Name,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign token with secret key
	tokenString, _ := token.SignedString([]byte("secret"))

	return tokenString
}

func HashPassword(password string) (string, error) {

	// Generate a salt
	salt := bcrypt.DefaultCost

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), salt)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
