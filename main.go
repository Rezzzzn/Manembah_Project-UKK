package main

import (
	"go-backend-basic/config"
	"go-backend-basic/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	// 🔥 CONNECT DATABASE
	config.ConnectDB()

	// 🔐 GOOGLE OAUTH CONFIG
	config.GoogleOAuthConfig = &oauth2.Config{
		ClientID:     "78018392363-v7fblhu0b6crsrt1f08esk9bkc799rvq.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	r := gin.Default()

	// 🌐 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://127.0.0.1:5500",
			"http://localhost:5500",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
		},
		AllowCredentials: true,
	}))

	// ROUTES
	routes.SetupRoutes(r)

	r.Run(":8080")
}
