	package handlers

import (
	"net/http"

	"go-backend-basic/config"
	"github.com/gin-gonic/gin"
)

func Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak valid"})
		return
	}

	var name, email string

	err := config.DB.QueryRow(
		"SELECT name, email FROM users WHERE id=$1",
		userID,
	).Scan(&name, &email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    userID,
		"name":  name,
		"email": email,
	})
}
