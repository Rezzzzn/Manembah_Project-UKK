package handlers

import (
	"database/sql"
	"net/http"

	"go-backend-basic/config"
	"go-backend-basic/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID int
	var name string
	var role string
	var hashedPassword string

	err := config.DB.QueryRow(
		"SELECT id, name, role, password FROM users WHERE email=$1",
		req.Email,
	).Scan(&userID, &name, &role, &hashedPassword)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email tidak ditemukan"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(req.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password salah"})
		return
	}

	// 🔐 GENERATE JWT
	token, err := utils.GenerateToken(userID, req.Email, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login berhasil",
		"token":   token,
		"user": gin.H{
			"id":    userID,
			"name":  name,
			"email": req.Email,
			"role":  role,
		},
	})
}

