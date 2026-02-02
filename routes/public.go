package routes

import (
	"go-backend-basic/handlers"

	"github.com/gin-gonic/gin"
)

func PublicRoutes(r *gin.Engine) {
	// AUTH
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// GOOGLE AUTH
	r.GET("/auth/google/login", handlers.GoogleLogin)
	r.GET("/auth/google/callback", handlers.GoogleCallback)

	// USERS (OPTIONAL)
	r.GET("/users", handlers.GetUsers)
	r.POST("/users", handlers.CreateUser)

	// PUBLIC READ
	r.GET("/accommodations", handlers.GetAccommodations)
}
