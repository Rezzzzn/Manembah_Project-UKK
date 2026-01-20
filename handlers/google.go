package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"go-backend-basic/config"
	"go-backend-basic/utils"

	"github.com/gin-gonic/gin"
)

type GoogleUser struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func GoogleLogin(c *gin.Context) {
	url := config.GoogleOAuthConfig.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")

	token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(400, gin.H{"error": "oauth exchange gagal"})
		return
	}

	client := config.GoogleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(400, gin.H{"error": "gagal ambil data google"})
		return
	}
	defer resp.Body.Close()

	var user GoogleUser
	json.NewDecoder(resp.Body).Decode(&user)

	// 🔍 CEK USER
	var id int
	var role string

	err = config.DB.QueryRow(
		"SELECT id, role FROM users WHERE email=$1",
		user.Email,
	).Scan(&id, &role)

	if err != nil {
		// AUTO REGISTER
		err = config.DB.QueryRow(
			"INSERT INTO users (name, email, role) VALUES ($1,$2,'user') RETURNING id, role",
			user.Name,
			user.Email,
		).Scan(&id, &role)
	}

	// 🔐 BUAT JWT
	tokenJWT, err := utils.GenerateToken(id, user.Email, role)
	if err != nil {
		c.JSON(500, gin.H{"error": "gagal generate token"})
		return
	}

	// 🔁 REDIRECT KE FRONTEND
	redirectURL := "http://127.0.0.1:5500/dashboard.html?token=" + tokenJWT
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

