package middleware

import (
	"big-devops-api/internal/services"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(authSvc services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
		}

		claims, err := authSvc.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token: " + err.Error()})
		}

		// Store claims in context
		c.Locals("user", claims)

		// Extract user_id
		var userID uint
		if uid, ok := claims["user_id"].(float64); ok {
			userID = uint(uid)
		} else if uid, ok := claims["user_id"].(uint); ok {
			userID = uid
		}
		c.Locals("user_id", userID)

		// Extract role
		role := ""
		if r, ok := claims["role"].(string); ok {
			role = r
		} else if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
			// Keycloak fallback
			if rRoles, ok := realmAccess["roles"].([]interface{}); ok {
				for _, r := range rRoles {
					roleStr := r.(string)
					// Prioritize admin or operator roles if found
					if roleStr == "admin" || roleStr == "operator" {
						role = roleStr
						break
					}
					if role == "" {
						role = roleStr
					}
				}
			}
		}
		c.Locals("role", role)

		return c.Next()
	}
}

func RoleChecker(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals("user").(jwt.MapClaims)
		
		var roles []string
		
		// Try Keycloak format
		if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
			if rRoles, ok := realmAccess["roles"].([]interface{}); ok {
				for _, role := range rRoles {
					roles = append(roles, role.(string))
				}
			}
		}
		
		// Try Local Auth format
		if role, ok := claims["role"].(string); ok {
			roles = append(roles, role)
		}

		hasRole := false
		for _, role := range roles {
			if role == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
		}

		return c.Next()
	}
}
