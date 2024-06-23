package middleware

import (
	"ballot-box/internal/utils/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AdminOnlyMiddleware() gin.HandlerFunc {
	return func (c *gin.Context) {

		userRole := c.GetString("role")

		if userRole != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin access required"})
            c.Abort()
            return
		}
		c.Next()
	}
}

func VoterOnlyMiddleware() gin.HandlerFunc {
	return func (c *gin.Context) {

		userRole := c.GetString("role")

		if userRole != "voter" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Voter access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {

        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

        token, err := auth.VerifyToken(tokenString)

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusBadGateway, gin.H{"error": "error getting token claims"})
            return
        }

        role, rok := claims["role"].(string)
        userID, uok := claims["id"].(string)

        if !uok || !rok {
            c.JSON(http.StatusBadGateway, gin.H{"error":  "error getting claims"})
            return
        }

        c.Set("userID", userID)
        c.Set("role", role)
        c.Next()
    }
}