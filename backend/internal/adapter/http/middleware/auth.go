package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"workshop-registration/backend/internal/adapter/storage/postgres/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.GetHeader("Authorization")

		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing authentication token",
			})
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")

		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			},
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithExpirationRequired(),
		)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		id, err := strconv.ParseUint(claims.Subject, 10, 64)

		if err != nil || id == 0 || id > uint64(^uint(0)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user identity",
			})
			return
		}

		var user models.UserModel

		if err := db.First(&user, uint(id)).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User account not available",
			})
			return
		}

		c.Set("userID", user.ID)
		c.Set("role", user.Role)

		c.Next()
	}
}

func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		currentRole := c.GetString("role")

		for _, role := range roles {
			if currentRole == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Access denied",
		})
	}
}
