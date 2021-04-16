package services

import (
	errors2 "errors"
	"fmt"
	"github.com/Ekod/highload-otus/utils/errors"
	"github.com/Ekod/highload-otus/utils/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"os"
	"time"
)

var (
	secretKey                                = os.Getenv("jwt_secret")
	SecurityService securityServiceInterface = &securityService{}
)

type securityService struct{}

type securityServiceInterface interface {
	HashPassword(password string) (string, *errors.RestErr)
	VerifyPassword(password string, hashedPassword string) *errors.RestErr
	GenerateToken(userId int64) (string, *errors.RestErr)
	ParseToken(accessToken string) (int64, *errors.RestErr)
	GetUserIdFromToken(c *gin.Context) (int64, *errors.RestErr)
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int64 `json:"userId"`
}

//HashPassword хэширует пароль
func (s *securityService) HashPassword(password string) (string, *errors.RestErr) {
	hp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.LogError("Services_Security_HashPassword - error while hashing the password", err)
		return "", errors.NewInternalServerError(fmt.Sprintf("error while hashing the password: %s", err.Error()))
	}

	return string(hp), nil
}

//VerifyPassword проверяет валидность приходящего пароля при логине
func (s *securityService) VerifyPassword(password string, hashedPassword string) *errors.RestErr {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		logger.LogError("Services_Security_VerifyPassword - error while verifying the password", err)
		return errors.NewBadRequestError(fmt.Sprintf("Email or password is invalid"))
	}

	return nil
}

//GenerateToken генерирует jwt-токен
func (s *securityService) GenerateToken(userId int64) (string, *errors.RestErr) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		logger.LogError("Services_Security_GenerateToken - error while generating token", err)
		return "", errors.NewInternalServerError(fmt.Sprintf("error while generating token"))
	}
	return tokenString, nil
}

//ParseToken проверяет валидность токена
func (s *securityService) ParseToken(accessToken string) (int64, *errors.RestErr) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors2.New("invalid signing method")
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, errors.NewBadRequestError("invalid signing method")
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.NewBadRequestError("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

//GetUserIdFromToken получает id пользователя после функции middlewares.CheckIdInToken
func (s *securityService) GetUserIdFromToken(c *gin.Context) (int64, *errors.RestErr) {
	id, ok := c.Get("userId")
	if !ok {
		return 0, errors.NewInternalServerError("user id not found")
	}

	idInt, ok := id.(int64)
	if !ok {
		return 0, errors.NewInternalServerError("user id is not a number")
	}

	return idInt, nil
}
