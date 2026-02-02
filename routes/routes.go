package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine) {
	PublicRoutes(r)
	ProtectedRoutes(r)
}
