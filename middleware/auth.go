package middleware

import (
	"project/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)


func VerifyToken(c *fiber.Ctx) error{
	authHeader:=c.Get("Authorization")

	parts:=strings.Split(authHeader," ")

	if len(parts)!=2 || parts[0]!= "Bearer"{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
	}  
	
	tokenStr := parts[1]
	token ,err :=jwt.Parse(tokenStr,utils.ExtractSecriteKey) 
	if err!=nil || !token.Valid{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	userId, ok:=claims["user_id"]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token data"})
	}

	
	role, ok:= claims["role"]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token data"})
	}

	validRoles := map[string]bool{
    "patient": true,
    "doctor":  true,	
	"admin":true,
	}

	roleStr, ok := role.(string) 
	if !ok || !validRoles[roleStr] {
    	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid role"})
	}

	c.Locals("user_id", userId)
	c.Locals("role", roleStr) 
	return c.Next() 
}



