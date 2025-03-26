package routes

import (
	"github.com/amirtalbi/examen_go/internal/api/handlers"
	"github.com/amirtalbi/examen_go/internal/api/middleware"
	"github.com/amirtalbi/examen_go/internal/config"
	"github.com/amirtalbi/examen_go/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, authService service.AuthService, userService service.UserService) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.LoggerMiddleware())

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	healthHandler := handlers.NewHealthHandler()

	apiGroup := router.Group("/" + cfg.APIPrefix)

	apiGroup.GET("/health", healthHandler.Check)

	authRoutes := apiGroup.Group("/")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/forgot-password", authHandler.ForgotPassword)
		authRoutes.POST("/reset-password", authHandler.ResetPassword)
		// Moved refresh endpoint outside of protected routes
		authRoutes.POST("/refresh", authHandler.RefreshToken)
	}

	protected := apiGroup.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		protected.POST("/logout", authHandler.Logout)
		protected.GET("/me", userHandler.GetProfile)
	}

	return router
}
