package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type MiddlewareVerify struct {
	DB     *gorm.DB
	Secret string
}

func SetUpMiddleware(db *gorm.DB, secret string) MiddlewareVerify {
	return MiddlewareVerify{
		DB:     db,
		Secret: secret,
	}
}

func (m *MiddlewareVerify) UserAuth(ctx *gin.Context) {
	token, err := getTokenFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	parsedToken, err := verifyToken(token, m.Secret)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		ctx.Abort()
		return
	}

	fmt.Println("Claims are:", claims)

	ctx.Set("user", claims)

	ctx.Next()
}

func getTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}
	parts := strings.Split(authHeader, " ")
	fmt.Printf("authHeader:%v\n", parts)
	if len(parts) < 2 || parts[0] != "Bearer" {
		fmt.Printf("parts[0] is:%v\n", len(parts[0]))
		return "", errors.New("invalid token")
	}
	return parts[1], nil
}

func verifyToken(tokenString string, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
