package routes

import (
	"go-backend-basic/handlers"
	"go-backend-basic/middleware"

	"github.com/gin-gonic/gin"
)

func ProtectedRoutes(r *gin.Engine) {
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/profile", handlers.Profile)
		auth.POST("/logout", handlers.Logout)

		AdminRoutes(auth)
	}
}
