package routes

import (
	"go-backend-basic/handlers"
	"go-backend-basic/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// ======================
	// PUBLIC ROUTES
	// ======================
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// Google Login
	r.GET("/auth/google/login", handlers.GoogleLogin)
	r.GET("/auth/google/callback", handlers.GoogleCallback)

	// Optional
	r.GET("/users", handlers.GetUsers)
	r.POST("/users", handlers.CreateUser)

	// ======================
	// PUBLIC READ
	// ======================
	r.GET("/accommodations", handlers.GetAccommodations)

	// ======================
	// PROTECTED ROUTES
	// ======================
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/profile", handlers.Profile)
		auth.POST("/logout", handlers.Logout)

		// ======================
		// ADMIN ONLY
		// ======================
		admin := auth.Group("/admin")
		admin.Use(middleware.RequireRole("admin"))
		{
			// DASHBOARD
			admin.GET("/dashboard", handlers.AdminDashboard)

			// CRUD ACCOMMODATIONS
			admin.POST("/accommodations", handlers.CreateAccommodation)
			admin.PUT("/accommodations/:id", handlers.UpdateAccommodation)
			admin.DELETE("/accommodations/:id", handlers.DeleteAccommodation)

			// MANAGE USERS
			admin.GET("/users", handlers.GetAllUsers)
			admin.GET("/users/:id", handlers.GetUserByID)
			admin.PUT("/users/:id", handlers.UpdateUser)
			admin.DELETE("/users/:id", handlers.DeleteUser)
		}
	}
}
