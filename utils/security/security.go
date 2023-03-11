package security

import (
	"fmt"
	"github.com/Ekod/highload-otus/utils/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"os"

	"time"
)

var secretKey = os.Getenv("jwt_secret")

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"userId"`
}

//HashPassword хэширует пароль
func HashPassword(password string) (string, error) {
	hp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] Services_Security_HashPassword - error while hashing the password: %s", err)

		return "", errors.NewInternalServerError("Server error", debugMessageError)
	}

	return string(hp), nil
}

//VerifyPassword проверяет валидность приходящего пароля при логине
func VerifyPassword(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] Services_Security_VerifyPassword - error while verifying the password: %s", err)

		return errors.NewBadRequestError("Email or password is invalid", debugMessageError)
	}

	return nil
}

//GenerateToken генерирует jwt-токен
func GenerateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] Services_Security_GenerateToken - error while generating token: %s", err)

		return "", errors.NewInternalServerError("Server error", debugMessageError)
	}
	return tokenString, nil
}

//ParseToken проверяет валидность токена
func ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			debugMessageError := "[ERROR] Services_Security_ParseToken - error while parsing token)"

			return nil, errors.NewInternalServerError("Server error", debugMessageError)
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] Services_Security_ParseToken - error while parsing token: %s", err)

		return 0, errors.NewBadRequestError("invalid signing method", debugMessageError)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		debugMessageError := fmt.Sprintf("[ERROR] Services_Security_ParseToken - error while parsing token: %s", err)

		return 0, errors.NewBadRequestError("token claims are not of the correct type", debugMessageError)
	}

	return claims.UserId, nil
}

//GetUserIdFromToken получает id пользователя после функции middlewares.CheckIdInToken
func GetUserIdFromToken(c *gin.Context) (int, error) {
	id, ok := c.Get("userId")
	if !ok {
		debugMessageError := "[ERROR] Services_Security_GetUserIdFromToken - error while getting id from token"

		return 0, errors.NewInternalServerError("user_service id not found", debugMessageError)
	}

	idInt, ok := id.(int)
	if !ok {
		debugMessageError := "[ERROR] Services_Security_GetUserIdFromToken - error while parsing token"

		return 0, errors.NewInternalServerError("user_service id is not a number", debugMessageError)
	}

	return idInt, nil
}
