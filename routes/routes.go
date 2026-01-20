package routes

import (
	"go-backend-basic/handlers"
	"go-backend-basic/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// PUBLIC ROUTES
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// OPTIONAL
	r.GET("/users", handlers.GetUsers)
	r.POST("/users", handlers.CreateUser)

	// 🔒 PROTECTED ROUTES
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/profile", handlers.Profile)
		auth.POST("/logout", handlers.Logout)

		// 🔐 ADMIN ONLY
		admin := auth.Group("/admin")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.GET("/dashboard", handlers.AdminDashboard)
		}
	}

	// login google

	r.GET("/auth/google/login", handlers.GoogleLogin)
	r.GET("/auth/google/callback", handlers.GoogleCallback)
}

