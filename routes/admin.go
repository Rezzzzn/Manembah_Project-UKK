package routes

import (
	"go-backend-basic/handlers"
	"go-backend-basic/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(auth *gin.RouterGroup) {
	admin := auth.Group("/admin")
	admin.Use(middleware.RequireRole("admin"))
	{
		// DASHBOARD
		admin.GET("/dashboard", handlers.AdminDashboard)

		// ACCOMMODATIONS CRUD
		admin.GET("/accommodations", handlers.GetAccommodations)
		admin.GET("/accommodations/:id", handlers.GetAccommodationByID)
		admin.POST("/accommodations", handlers.CreateAccommodation)
		admin.PUT("/accommodations/:id", handlers.UpdateAccommodation)
		admin.DELETE("/accommodations/:id", handlers.DeleteAccommodation)

		// USER MANAGEMENT
		admin.GET("/users", handlers.GetAllUsers)
		admin.GET("/users/:id", handlers.GetUserByID)
		admin.PUT("/users/:id", handlers.UpdateUser)
		admin.DELETE("/users/:id", handlers.DeleteUser)
	}
}
