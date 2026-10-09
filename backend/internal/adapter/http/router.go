package http

import (
	"net/http"

	"workshop-registration/backend/internal/adapter/http/api"
	"workshop-registration/backend/internal/adapter/http/middleware"
	"workshop-registration/backend/internal/bootstrap"

	"github.com/gin-gonic/gin"
)

func SetupRouter(app *bootstrap.App) *gin.Engine {
	router := gin.Default()

	authHandler := api.NewAuthHandler(app.AuthService)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Workshop Registration API running",
		})
	})

	router.POST("/api/auth/login", authHandler.Login)

	protected := router.Group("/api")
	protected.Use(
		middleware.AuthMiddleware(app.DB, app.Config.JWTSecret),
	)

	admin := protected.Group("")
	admin.Use(middleware.RequireRoles("admin"))
	admin.POST("/users", authHandler.CreateUser)

	return router
}
