package handlers

import (
	"database/sql"
	"net/http"

	"go-backend-basic/config"

	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

//
// =====================
// GET ALL USERS (ADMIN)
// =====================
//
func GetAllUsers(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, name, email, role
		FROM users
		ORDER BY id ASC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	users := make([]UserResponse, 0)

	for rows.Next() {
		var u UserResponse
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		users = append(users, u)
	}

	// cek error setelah loop
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

//
// =====================
// GET USER BY ID
// =====================
//
func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	var user UserResponse
	err := config.DB.QueryRow(`
		SELECT id, name, email, role
		FROM users
		WHERE id = $1
	`, id).Scan(&user.ID, &user.Name, &user.Email, &user.Role)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user tidak ditemukan",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

//
// =====================
// UPDATE USER
// =====================
//
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Name string `json:"name" binding:"required"`
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "input tidak valid",
		})
		return
	}

	// validasi role
	if input.Role != "admin" && input.Role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "role harus admin atau user",
		})
		return
	}

	result, err := config.DB.Exec(`
		UPDATE users
		SET name = $1, role = $2
		WHERE id = $3
	`, input.Name, input.Role, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user berhasil diupdate",
	})
}

//
// =====================
// DELETE USER
// =====================
//
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	result, err := config.DB.Exec(`
		DELETE FROM users
		WHERE id = $1
	`, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user berhasil dihapus",
	})
}
