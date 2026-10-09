package http

import (
	"net/http"

	"github.com/gin-contrib/cors"

	"workshop-registration/backend/internal/adapter/http/api"
	"workshop-registration/backend/internal/adapter/http/middleware"
	"workshop-registration/backend/internal/bootstrap"

	"github.com/gin-gonic/gin"
)

func SetupRouter(app *bootstrap.App) *gin.Engine {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
	}

	corsConfig.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Authorization",
	}

	corsConfig.AllowMethods = []string{
		"GET", "POST", "PUT", "OPTIONS",
	}

	router.Use(cors.New(corsConfig))

	authHandler := api.NewAuthHandler(app.AuthService)
	workshopHandler := api.NewWorkshopHandler(app.WorkshopService)
	registrationHandler := api.NewRegistrationHandler(app.RegistrationService)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	})

	router.POST("/api/auth/login", authHandler.Login)

	protected := router.Group("/api")
	protected.Use(
		middleware.AuthMiddleware(app.DB, app.Config.JWTSecret),
	)

	protected.GET("/auth/me", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"id":   c.GetUint("userID"),
			"role": c.GetString("role"),
		})
	})

	admin := protected.Group("")
	admin.Use(middleware.RequireRoles("admin"))
	admin.POST("/users", authHandler.CreateUser)

	manager := protected.Group("")
	manager.Use(middleware.RequireRoles("manager"))
	manager.POST("/workshops", workshopHandler.Create)
	manager.PUT("/workshops/:id", workshopHandler.Update)

	staff := protected.Group("")
	staff.Use(middleware.RequireRoles("manager", "staff"))

	staff.GET("/workshops", workshopHandler.GetAll)
	staff.GET("/workshops/:id", workshopHandler.GetByID)

	staff.POST(
		"/workshops/:id/registrations",
		registrationHandler.Create,
	)

	staff.GET(
		"/workshops/:id/registrations",
		registrationHandler.GetByWorkshop,
	)

	staff.POST(
		"/registrations/:id/cancel",
		registrationHandler.Cancel,
	)

	return router
}
