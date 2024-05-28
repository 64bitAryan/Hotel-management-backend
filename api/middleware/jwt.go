package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("Unauthorised")
	}
	if err := parseToken(token[0]); err != nil {
		return err
	}
	return nil
}

func parseToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("Unauthorised")
		}
		secret := os.Getenv("JWT_SECRET")
		fmt.Println("NEVER PRINT SECREAT: ", secret)
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT token: ", err)
		return fmt.Errorf("Unauthorised")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	}
	return fmt.Errorf("Unauthorised")
}
