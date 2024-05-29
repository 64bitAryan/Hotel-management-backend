package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("Unauthorised")
	}
	claims, err := validateToken(token[0])
	if err != nil {
		return err
	}
	// check token expiration
	expires := claims["expires"].(float64)
	if time.Now().Unix() > int64(expires) {
		return fmt.Errorf("token expired")
	}
	return c.Next()
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("Unauthorised")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT token: ", err)
		return nil, fmt.Errorf("Unauthorised")
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("Unauthorised")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Unauthorised")
	}
	return claims, nil
}
